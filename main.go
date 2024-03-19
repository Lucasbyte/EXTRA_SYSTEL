package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/text/encoding/charmap"
)

func main() {
	for {
		arquivo := "../itensmgv.txt"
		if _, err := os.Stat("../itensmgv.bak"); err == nil {
			arquivo = "../itensmgv.bak"
		}

		if _, err := os.Stat(arquivo); err == nil {
			time.Sleep(1 * time.Second)

			var dict_conserva, dict_fraciona, dict_aler, dict_forn map[string]string
			var info map[string]string

			if _, err := os.Stat("../conserva.bak"); err == nil {
				dict_conserva = conservaAnalyze("../conserva.bak")
			} else if _, err := os.Stat("../conserva.txt"); err == nil {
				dict_conserva = conservaAnalyze("../conserva.txt")
			}

			if _, err := os.Stat("../fraciona.bak"); err == nil {
				dict_fraciona = fracionaAnalyze("../fraciona.bak")
			} else if _, err := os.Stat("../fraciona.txt"); err == nil {
				dict_fraciona = fracionaAnalyze("../fraciona.txt")
			}

			if _, err := os.Stat("../campext1.bak"); err == nil {
				dict_aler = alergiaAnalyze("../campext1.bak")
			} else if _, err := os.Stat("../campext1.txt"); err == nil {
				dict_aler = alergiaAnalyze("../campext1.txt")
			}

			if _, err := os.Stat("../txforn.bak"); err == nil {
				dict_forn = fornAnalyze("../txforn.bak")
			} else if _, err := os.Stat("../txforn.txt"); err == nil {
				dict_forn = fornAnalyze("../txforn.txt")
			}

			if _, err := os.Stat("../txinfo.bak"); err == nil {
				info = infoAnalyze("../txinfo.bak")
			} else if _, err := os.Stat("../txinfo.txt"); err == nil {
				info = infoAnalyze("../txinfo.txt")
			}
			arquivoInfonutri := "../infnutri.txt"
			if _, err := os.Stat("../infnutri.bak"); err == nil {
				arquivoInfonutri = "../infnutri.bak"
			}

			infonutriFile, err := os.Open(arquivoInfonutri)
			if err != nil {
				fmt.Println(err)
			}
			defer infonutriFile.Close()

			mgv7File, err := os.Open("../itens.TXT")
			if err != nil {
				fmt.Println(err)
			}
			defer mgv7File.Close()

			// Create output files
			nutriFile, err := os.Create("nutriSystel.TXT")
			if err != nil {
				fmt.Println(err)
			}
			defer nutriFile.Close()

			systelFile, err := os.Create("itensSystel.TXT")
			if err != nil {
				fmt.Println(err)
			}
			defer systelFile.Close()

			// Initialize arrays
			var codPluArray []string
			var codNutriArray []string
			var codNutriMGVArray []string
			var receitaArray []string

			// Read and process mgv7File
			scanner := bufio.NewScanner(mgv7File)
			for scanner.Scan() {
				line := scanner.Text()

				codPlu := line[3:9]
				codPluArray = append(codPluArray, codPlu)

				codNutriMGV := line[78:84]
				codNutriMGVArray = append(codNutriMGVArray, codNutriMGV)

				codReceita := line[68:74]
				receitaArray = append(receitaArray, codReceita)

				textModified := line[0:43] + strings.Repeat(" ", 25) + codPlu + line[74:150] +
					"000000|01|                                                                      0000000000000000000000000||0||0000000000000000000000"
				fmt.Fprintln(systelFile, textModified)
			}

			// Read and process infonutriFile
			scanner = bufio.NewScanner(infonutriFile)
			for scanner.Scan() {
				line := scanner.Text()

				codNutri := line[1:7]
				boo := true
				if len(line) > 106 {
					boo = line[7:110] != "000000000000000000000000000000000000000000|000000000000000000000000000000000000000000000000000000000000"
				}
				if len(line) < 50 {
					line = line[0:49] + "\n"
					line += "|" + strings.Repeat("0", 3)
					porcao := line[7:11]
					if parseInt(porcao) <= 0 {
						porcao = "0100"
					}
					line += porcao + "0" + line[12:26]
					line += strings.Repeat("0", 6) + line[26:50]
					line += strings.Repeat("0", 9)
				} else {
					fmt.Println(line)
					line = strings.ReplaceAll(line, "|", "0000|")
				}

				if parseInt(line[61:63]) > 28 {
					line = line[:61] + "16" + line[63:]
				}

				if containsTo(codNutriMGVArray, codNutri) && !containsTo(codNutriArray, codNutri) && boo {
					codNutriArray = append(codNutriArray, codNutri)
					fmt.Println("oi", line)
					fmt.Fprintln(nutriFile, line)
				} else {
					fmt.Println(codNutriMGVArray[0], codNutriArray[0], codNutri, boo)
				}
			}

			fmt.Println("Data analysis completed successfully.")

			itensAnalyze(arquivo)

			// arquivoMGV7, err := os.Open("../itens.TXT")
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }
			// defer arquivoMGV7.Close()

			// nutri, err := os.Create("nutriSystel.TXT")
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }
			// defer nutri.Close()

			// arquivoSystel, err := os.Create("itensSystel.TXT")
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }
			// defer arquivoSystel.Close()

			// scanner := bufio.NewScanner(arquivoMGV7)
			// for scanner.Scan() {
			// 	linha := scanner.Text()
			// 	codPlu := linha[3:9]
			// 	codNutriMGV := linha[79:84]
			// 	//codReceita := linha[68:74]

			// 	textoModificado := linha[0:43] + strings.Repeat(" ", 25) + codPlu + linha[74:150] + "000000|01|                                                                      0000000000000000000000000||0||0000000000000000000000"
			// 	arquivoSystel.WriteString(textoModificado + "\n")

			// 	if _, err := strconv.Atoi(codNutriMGV); err == nil {
			// 		_, err := strconv.Atoi(codNutriMGV)
			// 		if err == nil {
			// 			nutri.WriteString(linha[0:49] + "|" + strings.Repeat("0", 3) + codNutriMGV + strings.Repeat("0", 6) + linha[26:50] + strings.Repeat("0", 9) + "\n")
			// 		}
			// 	}

			if err := scanner.Err(); err != nil {
				fmt.Println(err)
				return
			}

			arr_ip := comunicabal()
			if !containsTo(arr_ip, "localhost") {
				arr_ip = append(arr_ip, "localhost")
			}

			if len(arr_ip) > 0 {
				txt := ""
				for _, ip := range arr_ip {
					ip_db := strings.TrimSpace(ip)
					err = infoSystelWriter("itensSystel.TXT", info, dict_forn, dict_aler, dict_fraciona, dict_conserva, ip_db)
					if err != nil {
						logToFile("log-erro-conexao-go.txt", fmt.Sprintf("erro ao importar para: %s erro: %s\n ", ip, err.Error()))
					}
					txt += fmt.Sprintf("%s ", ip)
					time.Sleep(5 * time.Second)

				}
				logToFile("log.txt", fmt.Sprintf("importou corretamente: %s\n", txt))
			}

			fmt.Println("pronto")
			time.Sleep(25 * time.Second)
		}

		time.Sleep(20 * time.Second)
	}
}

func conservaAnalyze(arq string) map[string]string {
	conservaTxt, err := os.Open(arq)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer conservaTxt.Close()

	dictConserva := make(map[string]string)
	scanner := bufio.NewScanner(conservaTxt)
	for scanner.Scan() {
		line := scanner.Text()
		key := line[0:4]
		value := caracterRemove(line[104:])
		if key != "0003" {
			value = caracterRemove(line[4:40])
		}
		dictConserva[key] = value
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}

	return dictConserva
}

func fracionaAnalyze(arq string) map[string]string {
	fracionaTxt, err := os.Open(arq)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer fracionaTxt.Close()

	dictFraciona := make(map[string]string)
	scanner := bufio.NewScanner(fracionaTxt)
	for scanner.Scan() {
		line := scanner.Text()
		dictFraciona[line[0:4]] = caracterRemove(line[104:])
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}

	return dictFraciona
}

func alergiaAnalyze(arq string) map[string]string {
	campTxt, err := os.Open(arq)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer campTxt.Close()

	dictAler := make(map[string]string)
	scanner := bufio.NewScanner(campTxt)
	for scanner.Scan() {
		line := scanner.Text()
		dictAler[line[0:4]] = caracterRemove(line[104:])
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}

	return dictAler

}

func fornAnalyze(arq string) map[string]string {
	fornTxt, err := os.Open(arq)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer fornTxt.Close()

	dictFornecedor := make(map[string]string)
	scanner := bufio.NewScanner(fornTxt)
	for scanner.Scan() {
		line := scanner.Text()
		dictFornecedor[line[0:4]] = caracterRemove(line[104:217])
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}

	return dictFornecedor
}

func infoAnalyze(arq string) map[string]string {
	infoTxt, err := os.Open(arq)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer infoTxt.Close()

	dictInfo := make(map[string]string)
	scanner := bufio.NewScanner(infoTxt)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("linha:       ", line)
		if len(line) > 100 {
			dictInfo[line[0:6]] = caracterRemove(line[106:])
		} else if len(line) > 6 {
			dictInfo[line[0:6]] = caracterRemove(line)
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}

	return dictInfo
}

func infoSystelWriter(arq string, d_info, d_forn, d_aler, d_fra, d_con map[string]string, ip string) error {
	item, err := os.Open(arq)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer item.Close()

	passw, user := "1234", "user"
	if ip == "localhost" {
		passw, user = "Systel#4316", "systel"
	}

	db_acess := fmt.Sprintf("user=%s password=%s host=%s dbname=cuora sslmode=disable", user, passw, ip)
	db, err := sql.Open("postgres", db_acess)
	if err != nil {
		fmt.Println("ERRO: ", err)
		return err
	}
	defer db.Close()

	scanner := bufio.NewScanner(item)
	for scanner.Scan() {
		line := scanner.Text()
		lote := line[90:102]
		codPlu := line[3:9]
		codInfo := line[68:74]
		codAler := line[126:130]
		codForn := line[86:90]
		codFrac := line[122:126]
		codCons := line[134:138]

		info := d_info[codInfo]
		aler := d_aler[codAler]
		forn := d_forn[codForn]
		cons := d_con[codCons]
		frac := d_fra[codFrac]

		if len(forn) < 158 && forn != "" {
			forn += strings.Repeat(" ", 158-len(forn))
		}

		fornFrac := forn + " " + frac

		loteInt, _ := strconv.Atoi(lote)

		if loteInt > 0 {
			enviarInf(strconv.Itoa(loteInt), codPlu, db, "lot")
		}
		if info != "" {
			enviarInf(info, codPlu, db, "extra_field1")
		}
		if aler != "" {
			enviarInf(aler, codPlu, db, "extra_field2")
		}
		if cons != "" {
			enviarInf(cons, codPlu, db, "preservation_info")
		}
		if forn != "" && frac != "" {
			enviarInf(fornFrac, codPlu, db, "ingredients")
		}
	}

	fmt.Println(EncodeToUTF8(ip))
	return nil
}

func itensAnalyze(arq string) {
	arquivoItensMgv, err := os.Open(arq)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer arquivoItensMgv.Close()

	arquivoMGV7, err := os.Create("../itens.TXT")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer arquivoMGV7.Close()

	scanner := bufio.NewScanner(arquivoItensMgv)
	for scanner.Scan() {
		linha := scanner.Text()
		textoModificado := linha
		if len(linha) < 60 {
			textoModificado = strings.ReplaceAll(linha, "\n", " ")
		}
		arquivoMGV7.WriteString(textoModificado + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}
}

func EncodeToUTF8(str string) (string, error) {
	encoder := charmap.Windows1252.NewEncoder()
	utf8Str, err := encoder.String(str)
	if err != nil {
		return "", err
	}
	return utf8Str, nil
}

func enviarInf(t, p string, db *sql.DB, campo string) {
	plu := p
	tara := t
	tara, err := EncodeToUTF8(tara)
	if strings.TrimSpace(strings.ReplaceAll(tara, "\n", "")) != "" {
		if err != nil {
			fmt.Println("ERRO: ", err, "item: ", tara)
		} else {
			fmt.Println("")
		}
		comando := fmt.Sprintf("UPDATE product set %s = '%s' WHERE product_id = %s", campo, tara, plu)
		_, err = db.Exec(comando)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func filePortsEx() []string {
	arquivo, err := os.Open("porta.txt")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer arquivo.Close()

	var arr []string
	scanner := bufio.NewScanner(arquivo)
	for scanner.Scan() {
		ip := strings.TrimSpace(scanner.Text())
		if ip != "" && ip != "-1" {
			arr = append(arr, ip)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}

	return arr
}

func setorWrite() {
	arquivo := "../setor.txt"
	if _, err := os.Stat(arquivo); os.IsNotExist(err) {
		setorFile, err := os.Create(arquivo)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer setorFile.Close()

		for i := 1; i <= 12; i++ {
			setor := strconv.Itoa(i)
			if len(setor) < 2 {
				setor = "0" + setor
			}
			nome := ""
			switch setor {
			case "01":
				nome = "GERAL"
			case "02":
				nome = "HORTIFRUTI"
			case "03":
				nome = "PADARIA"
			case "04":
				nome = "AÇOGUE"
			case "05":
				nome = "FRIOS"
			case "06":
				nome = "PREPARAÇÃO"
			case "07":
				nome = "NOBRE"
			case "08":
				nome = "HORTIFRUTI"
			case "09":
				nome = "GERAL 09"
			case "10":
				nome = "GERAL 10"
			case "11":
				nome = "GERAL 11"
			case "12":
				continue
			}

			setorFile.WriteString(setor + " " + nome + "\n")
		}
	}
}

func comunicabal() []string {
	arquivo, err := os.Open("porta.txt")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer arquivo.Close()

	var arr []string
	scanner := bufio.NewScanner(arquivo)
	for scanner.Scan() {
		ip := strings.TrimSpace(scanner.Text())
		if ip != "" && ip != "-1" {
			arr = append(arr, ip)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}

	return arr
}

func caracterRemove(txt string) string {
	texto := strings.ReplaceAll(txt, "\n", "")
	texto = strings.ReplaceAll(texto, "    ", "")
	texto = strings.ReplaceAll(texto, ", ", ",")
	texto = strings.ReplaceAll(texto, " (", "(")
	return texto
}

func containsTo(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func logToFile(file, text string) {
	log, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer log.Close()

	if _, err := log.WriteString(text); err != nil {
		fmt.Println(err)
	}
}

// Helper function to convert string to integer
func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println(err)
	}
	return i
}

// Helper function to check if an integer is in an array
func contains(target int, arr []int) bool {
	for _, val := range arr {
		if val == target {
			return true
		}
	}
	return false
}
