import os
import time
import os.path
import postgresql



def caracter_remove(txt):
    texto = txt.replace("\n", "").replace("    ", "").replace(", ", ",").replace(" (", "(")
    return texto

def conserva_analyze(arq):
    conserva_txt = open(arq, 'r')
    dict_conserva = {}
    for line in conserva_txt:
        #dict_conserva[line[0:4]] = (line[4:194].replace('\n', '')).replace("    ", "")
        dict_conserva[line[0:4]] =  caracter_remove(line[104:194]) if line[0:4] != '0003' else caracter_remove(line[4:40])

    conserva_txt.close()
    return dict_conserva

def fraciona_analyze(arq):
    fraciona_txt = open(arq, 'r')
    dict_fraciona = {}
    for line in fraciona_txt:
        #dict_fraciona[line[0:4]] = (line[104:217].replace('\n', '')).replace("    ", "")
        dict_fraciona[line[0:4]] = caracter_remove(line[104:])
    fraciona_txt.close()
    return dict_fraciona

def alergia_analyze(arq):
    camp_txt = open(arq, 'r')
    dict_aler = {}
    for line in camp_txt:
        #dict_aler[line[0:4]] = caracter_remove(line[104:272].replace('\n', '')).replace("    ", "")
        dict_aler[line[0:4]] = caracter_remove(line[104:])
    camp_txt.close()
    return dict_aler

def forn_analyze(arq):
    forn_txt = open(arq, 'r')
    dict_fornecedor = {}
    for line in forn_txt:

        dict_fornecedor[line[0:4]] = caracter_remove(line[104:217])
    forn_txt.close()
    return dict_fornecedor


def info_analyze(arq):
    info_txt = open(arq, 'r')
    dict_info = {}
    for line in info_txt:

        dict_info[line[0:6]] = caracter_remove(line[106:])
    info_txt.close()
    return dict_info


def infoSystel_writer(arq, d_info, d_forn, d_aler, d_fra, d_con, ip):
    item = open(arq, 'r')

    if True:
        try: pass
        except: print('erro')
        passw = '1234'
        user = 'user'
        if ip == 'localhost':
            passw = 'postgres'
            user = 'postgres'
        
        if True:
            db_acess = f'pq://{user}:{passw}@{ip}:5432/cuora'
            db = postgresql.open(db_acess)
            for line in item:
                lote = line[90:102]
                cod_plu = line[3:9]
                cod_info = line[68:74]
                cod_aler = line[126:130]
                cod_forn = line[86:90]
                cod_frac = line[122:126]
                cod_cons = line[134:138]
                espaco =  (' '*1)
                
                info = d_info[cod_info]  if cod_info in d_info else " "
                aler = d_aler[cod_aler] if cod_aler in d_aler else " "
                forn = d_forn[cod_forn] if cod_forn in d_forn else " "
                cons = d_con[cod_cons] if cod_cons in d_con else " "
                frac = d_fra[cod_frac] if cod_frac in d_fra else " "

                if len(frac) < 158 and frac != " ":
                    frac = frac + ' '*(158- len(frac))

                
                forn_frac = frac + ' ' + forn

                lote = int(lote)

                if lote > 0:
                    enviar_inf(str(lote), cod_plu, db, 'lot')
                if not(info ==  " "):
                    enviar_inf(d_info[cod_info], cod_plu, db, 'extra_field1')
                if not(aler ==  " "):
                    enviar_inf(d_aler[cod_aler], cod_plu, db, 'extra_field2')
                if not(cons ==  " "):
                    enviar_inf(d_con[cod_cons], cod_plu, db, 'preservation_info')
                if not(forn ==  " " and frac ==  " "):
                    enviar_inf(forn_frac, cod_plu, db, 'ingredients')

            db.close()
            print(ip)
        

    print('hm')
    item.close()


def itens_analize(arq):
    arquivoItensMgv = open(arq, 'r')
    arquivoMGV7 = open('../itens.TXT', 'w')
    for l in arquivoItensMgv:
        if len(l) < 60:
            textoModificado = l.replace("\n", " ")

        else: 
            textoModificado = l

        arquivoMGV7.write(textoModificado)

    arquivoItensMgv.close()
    arquivoMGV7.close()


#DB

def enviar_inf(t, p, db, campo):
    plu = p
    tara = t
    comando = f"""UPDATE product
               set {campo} = '{tara}'
               WHERE product_id = {plu}"""
    db.execute(comando)


def file_ports_ex():
    try:
        try:
            teste = open('porta.txt', 'r')
            teste.close()
            return open('porta.txt', 'r')
        except:
            try:
                teste2 = open('porta.txt', 'r')
                teste2.close()
                return open('porta.txt', 'r')
            except:
                teste3 = open('porta.txt', 'w')
                teste3.close()
                return open('porta.txt', 'r')
    except: print('erro ao ler arquivo de portas')





def setorWrite(array_setor):
    arquivo = '../setor.txt'
    if not (os.path.isfile(arquivo)):
        with open(arquivo, 'w') as setor_file:
            for i in range(1, 13):
                setor = str(i)
                if len(setor) < 2: setor = '0' + setor
                nome = ""
                match setor:
                    case '01':
                        nome = "GERAL"
                    case '02':
                        nome = "HORTIFRUTI"
                    case '03':
                        nome = "PADARIA"
                    case '04':
                        nome = "AÇOGUE"
                    case '05':
                        nome = "FRIOS"
                    case '06':
                        nome = "PREPARAÇÃO"
                    case '07':
                        nome = "NOBRE"
                    case '08':
                        nome = "HORTIFRUTI"
                    case '09':
                        nome = "GERAL 09"
                    case '10':
                        nome = "GERAL 10"
                    case '11':
                        nome = "GERAL 11"
                    case '12':
                        pass
                    





def comunicabal():
    array = []
    porta = file_ports_ex()
    for line in porta:
        ip = str(line)
        ip = ip.replace('\n', '')
        if ip != '' and ip != '-1':
            array.append(ip)
    porta.close()
    return array
                
                




while(True):
    
    arquivo = '../itensmgv.txt' if not (os.path.isfile('../itensmgv.bak')) else '../itensmgv.bak'
    
    if(os.path.isfile(arquivo)):
        time.sleep(1)

        
        

        if(os.path.isfile('../conserva.bak')):
            dict_conserva = conserva_analyze('../conserva.bak')
        
        elif(os.path.isfile('../conserva.txt')):
             dict_conserva = conserva_analyze('../conserva.txt')
        else: dict_conserva = ''

        if(os.path.isfile('../fraciona.bak')):
            dict_fraciona = fraciona_analyze('../fraciona.bak')
        elif(os.path.isfile('../fraciona.txt')):
            dict_fraciona = fraciona_analyze('../fraciona.txt')
        else: dict_fraciona = ''

        if(os.path.isfile('../campext1.bak')):
            dict_aler = alergia_analyze('../campext1.bak')
        elif(os.path.isfile('../campext1.txt')):
            dict_aler = alergia_analyze('../campext1.txt')
        else: dict_aler = ''

        if(os.path.isfile('../txforn.bak')):
            dict_forn =  forn_analyze('../txforn.bak')
        elif(os.path.isfile('../txforn.txt')):
            dict_forn =  forn_analyze('../txforn.txt')
        else:
            dict_forn = '' 
        
        if(os.path.isfile('../txinfo.bak')):
            info = info_analyze('../txinfo.bak')
        elif(os.path.isfile('../txinfo.txt')):
            info = info_analyze('../txinfo.txt')


        else:
            info = ""
            pass
        if(os.path.isfile('../infnutri.bak')):
            arquivoInfonutri = open('../infnutri.bak', 'r')
        elif(os.path.isfile('../infnutri.bak')):
            arquivoInfonutri = open('../infnutri.txt', 'r')
        else:
            arquivoInfonutri = open('../infnutri.bak', 'w')
        
        itens_analize(arquivo)


        arquivoMGV7 = open('../itens.TXT', 'r')

        nutri = open('nutriSystel.TXT', 'w')
        codPlu_array = []
        codNutri_array = []
        codNutriMGVARRAY = []
        receita_array = []
        setor_array = []
        arquivoSystel = open('itensSystel.TXT', 'w')
        for linha in arquivoMGV7:
        
            

            codPlu = linha[3:9]
            
            codPlu_array.append(codPlu)
            codNutriMGV = linha[79:84]
            codReceita = linha[68:74]
            receita_array.append(codReceita)

            codNutriMGVARRAY.append(int(codNutriMGV))
            textoModificado = linha[0:43] + (' ')*25 + codPlu + linha[74:150] +  "000000|01|                                                                      0000000000000000000000000||0||0000000000000000000000"


            arquivoSystel.write(textoModificado+"\n")

        arquivoSystel.close()



        
        array_cod_nutri = []


        


        for linha in arquivoInfonutri:
            codNutri = int(linha[1:7])

            boo = linha[7:110] != '000000000000000000000000000000000000000000|000000000000000000000000000000000000000000000000000000000000'
            #print(boo)
            if len(linha) < 50:
                #print(linha[7:11])
                linha = linha[0:49].replace('\n', '') 
                linha = linha + '|' + ('0'*3)
                porcao = linha[7:11] if int(linha[7:11]) > 0 else '0100'
                #print(porcao, int(linha[7:11]))
                linha = linha + porcao
                linha = linha + '0' + linha[12:26] 
                linha = linha + '0'*6 + linha[26:50].replace('\n','') 
                linha = linha + '0'*9 + '\n'
            else:
                linha = linha.replace('|', '0000|')
            if int(linha[61:63]) > 28:
                #print(linha[61:63])
                linha = linha[:61] + '16' + linha[63:]
            if codNutri in codNutriMGVARRAY and not(codNutri in array_cod_nutri) and boo:
                array_cod_nutri.append(codNutri)
                #if codNutri == 5:
                    #print("teste")
                #print(codNutri, codNutriMGVARRAY)
                nutri.write(linha)
                codNutri_array.append(codNutri)
            #print("1")
        arquivoMGV7.close()
        arquivoInfonutri.close()
        nutri.close()
        arquivoSystel.close()
        arr_ip = comunicabal()
        if not 'localhost' in arr_ip:
            arr_ip.append('localhost')
        if len(arr_ip) > 0:
            txt = ''
            for ip in arr_ip:
                #ip_db = ip.replace('\n', '')
                #infoSystel_writer('itensSystel.TXT', info, dict_forn, dict_aler, dict_fraciona, dict_conserva, ip_db)
                try:
                    ip_db = ip.replace('\n', '')
                    infoSystel_writer('itensSystel.TXT', info, dict_forn, dict_aler, dict_fraciona, dict_conserva, ip_db)
                except:
                    with open('log-erro-conexao.txt', 'a') as log:
                        log.write(f'erro ao importar para: {ip}\n')
                else:
                    txt += f'{ip} '
                    time.sleep(5)
                finally:
                    pass
            with open('log.txt', 'a') as log:
                log.write(f'importou corretamente: {txt}\n')
            
        #print(len(codPlu_array))
        #print(len(codNutri_array))
        #os.remove("../arqsok.bak")
        
        
        print('pronto')
        time.sleep(25)
                

    time.sleep(20)
    
