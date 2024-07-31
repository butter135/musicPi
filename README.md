# musicPi
SSH接続できるraspberrypiをyoutube(音楽)再生サーバー？にするためのプログラムです。
## 必要なもの
+ mpv : 動画再生ソフト
+ youtube-dl : mpvでyoutubeを再生するためのプラグイン
+ YoutubeDataAPIのキー : 環境変数"YTKEY"として登録しておく
## 準備
1. 名前をmusicPiにしてユーザーを作成
2. raspberrypiをSSH接続できるようにしておく
3. musicPiディレクトリに隠しディレクトリ .queue を作製
4. add.go, play.go, search.goをそれぞれビルドして、binディレクトリに置いておく。
## 使い方
1. musicPiにSSHで接続
2. search "検索キーワード" もしくは add "動画のURL" を使って動画をキューに追加する
3. play & で再生