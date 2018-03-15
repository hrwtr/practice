
## Get login cookie
```
curl  
    -X POST                     # HTTPメソッドの指定（-X）
    -k                          # SSLのエラーを無視して処理を継続（-k）
                                # （サーバー側証明書が不正、クライアント側のルート証明書が不正など） 
    -c ${source_cookies_path}   # クッキーの送信（-b）、保存（-c）  
    -F 'username='$USER         # フォームからファイルのアップロード（-Fもしくは--form）
    -F 'password'=$PASS 
    -v                          # 詳細をログ出力（-vもしくは--verbose）
    --url　${HUE_AUTH_URL}      #    
```
## Get dto list
```
curl 
    -sS                         # 進捗やエラーを表示しない（-sもしくは--silent）
                                # エラーだけは表示 (-S)	
    -X GET                      # HTTPメソッドの指定（-X）
    -k                          # SSLのエラーを無視して処理を継続（-k）
    -b ${source_cookies_path}   # クッキーの送信（-b）、保存（-c）  
    --url 
    -J                          # Content-Dispositionを参照し、＜ファイル名＞で指定された名前でファイルを取得する
    --output "${all_dto_list}"  # ダウンロードしたデータをファイルとして保存(-o --output)
```

## Export Json
```
curl 
    -k                          # SSLのエラーを無視して処理を継続（-k）
    -b ${source_cookies_path}   # クッキーの送信（-b）、保存（-c）  
    --url ${download_url}       
    -J                          # Content-Dispositionを参照し、＜ファイル名＞で指定された名前でファイルを取得する
    -O                          # ダウンロードしたデータをファイルとして保存(-o --output)
```
## Import Json
```
curl  
    -X POST                     # HTTPメソッドの指定（-X）
    -k                          # SSLのエラーを無視して処理を継続（-k）
    -b ${source_cookies_path}   # クッキーの送信（-b）、保存（-c）  
    -F                          # フォームからファイルのアップロード（-Fもしくは--form）
    --url ${upload_url}
```