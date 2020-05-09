# Review Bot

## 導入手順

1. ブラウザからSlackにログインする
![step1](img/step1.png)
![step2](img/step2.png)

2. https://api.slack.com/apps
へアクセスし、Create New Appをクリックする
![step3](img/step3.png)

3. 以下のように設定し、Create Appをクリックする  
**App Name** : review-bot  
**Development Slack Workspace** : 導入したいworkspace
![step4](img/step4.png)


4. 以下のような画面に遷移するので、Event Subscriptionsをクリックする
![step5](img/step5.png)


5. 以下のような設定にしてSaveする  
**Request URL** : このBotを動かしているサーバーのホスト/events-endpoint  
**Subscribe to bot events** : app_mention
![step6](img/step6.png)

6. OAuth & Permissionsのページを開く
![step7](img/step7.png)



7. 同ページのScopes > Bot Token Scopes
を以下のように設定
![step8](img/step8.png)


8. 同画面上部に戻って Install App to Workspace をクリックする
![step9](img/step9.png)


9. Allow をクリックする
![step10](img/step10.png)

10. OAuth Token が生成されるので控えておく
![step11](img/step11.png)

11. Basic Information > App Credentials の Vertification Token を控えておく
![step12](img/step12.png)

12. 10,11で控えた内容をこのBotを動かしているサーバーの環境変数に設定する **（複数の環境で運用する場合は:で区切って順番に設定する）**

```
export OAUTH_TOKENS=手順10で控えた値
export CLIENT_TOKENS=手順11で控えた値
```