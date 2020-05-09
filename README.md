# Review Bot

## 導入手順

1. ブラウザからSlackにログインする
![step1](https://user-images.githubusercontent.com/16130443/81471740-44e17a80-922e-11ea-9133-5b6c8fce7fb1.png)
![step2](https://user-images.githubusercontent.com/16130443/81471741-46ab3e00-922e-11ea-8d84-5bb661b84058.png)

2. https://api.slack.com/apps
へアクセスし、Create New Appをクリックする
![step3](https://user-images.githubusercontent.com/16130443/81471742-47dc6b00-922e-11ea-8c18-bb6378397272.png)

3. 以下のように設定し、Create Appをクリックする  
**App Name** : review-bot  
**Development Slack Workspace** : 導入したいworkspace
![step4](https://user-images.githubusercontent.com/16130443/81471743-48750180-922e-11ea-913c-b0ddbeb43cc6.png)


4. 以下のような画面に遷移するので、Event Subscriptionsをクリックする
![step5](https://user-images.githubusercontent.com/16130443/81471768-74908280-922e-11ea-8b88-db09d8d5f7e3.png)


5. 以下のような設定にしてSaveする  
**Request URL** : このBotを動かしているサーバーのホスト/events-endpoint  
**Subscribe to bot events** : app_mention
![step6](https://user-images.githubusercontent.com/16130443/81471769-778b7300-922e-11ea-8626-c7d7ad81cdac.png)

6. OAuth & Permissionsのページを開く
![step7](https://user-images.githubusercontent.com/16130443/81471771-78240980-922e-11ea-86fe-8651b48a4b9a.png)



7. 同ページのScopes > Bot Token Scopes
を以下のように設定
![step8](https://user-images.githubusercontent.com/16130443/81471773-78bca000-922e-11ea-9429-381844f4f92b.png)


8. 同画面上部に戻って Install App to Workspace をクリックする
![step9](https://user-images.githubusercontent.com/16130443/81471797-b4576a00-922e-11ea-8bbd-e0fc43d3fb71.png)


9. Allow をクリックする  
![step10](https://user-images.githubusercontent.com/16130443/81471799-b7525a80-922e-11ea-8306-82bff06f545c.png)

10. OAuth Token が生成されるので控えておく
![step11](https://user-images.githubusercontent.com/16130443/81471801-b7eaf100-922e-11ea-9e91-530193543dd4.png)

11. Basic Information > App Credentials の Vertification Token を控えておく
![step12](https://user-images.githubusercontent.com/16130443/81471803-b91c1e00-922e-11ea-8a7e-56d8a6f4eb12.png)

12. 10,11で控えた内容をこのBotを動かしているサーバーの環境変数に設定する **（複数の環境で運用する場合は:で区切って順番に設定する）**

```
export OAUTH_TOKENS=手順10で控えた値
export CLIENT_TOKENS=手順11で控えた値
```
