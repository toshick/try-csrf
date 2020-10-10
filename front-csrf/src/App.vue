<template>
  <main>
    <header>
      <h1>CSRFのための検証だよ</h1>
    </header>

    <div class="try">
      <div class="container">
        <section>
          <h2>前提知識</h2>
          <ul>
            <li>Cookieはポート番号では分離されないため、http://localhostもhttp://localhost:8888も同一ドメインと判断されるため注意が必要</li>
            <li>Cookieは「Top Level Navigationの場合ドメインをまたぐリクエストにも自動で付与される」という挙動（つまりページ遷移のこと）</li>
            <li>同じドメイン内でしか付与されなければ、例えば SiteA から SiteB に遷移した時に、必ずBのログイン画面が出てしまう（Cookieを利用したBの認証が認識されない）</li>
            <li>この実験はvue開発サーバ「http://127.0.0.1:8080/」でこのページを開き、goで立ち上げたhttpサーバ（{{ urlTop }}）へのリクエストを行う</li>
          </ul>
        </section>
        <!-- aリンク -->
        <section>
          <h2>
            <span class="ok">OK</span>
            <span>{{ urlUsers }}</span> へ移動
          </h2>
          <a :href="urlUsers">{{ urlUsers }}</a>

          <ul class="memo">
            <li>ページ遷移なのでcookie送信される</li>
            <li>cookieにSameSite:Laxなし → cookie送信される</li>
            <li>cookieにSameSite:Lax → cookie送信される</li>
            <li>Chrom、FireFox、Safariともに同じ</li>
          </ul>
        </section>
        <!-- formからPOST -->
        <section>
          <h2>
            <span class="warn">ブラウザにより違う</span>
            <span>{{ urlUsers }}</span> へForm Post
          </h2>
          <form :action="urlUsers" method="POST">
            <input name="myvalue" value="ispost" />
            <button type="submit">formからPOST</button>
          </form>
          <ul class="memo">
            <li>（Chrome）cookie送信されない</li>
            <li>（FireFox、Safari）cookie送信される</li>
            <li>cookieにSameSite:Laxなし → cookie送信される</li>
            <li>cookieにSameSite:Lax → cookie送信されない</li>
          </ul>
        </section>
        <!-- iframe -->
        <section>
          <h2>
            <span class="warn">ブラウザにより違う</span>
            <span>{{ urlUsers }}</span> へiframeでForm Post
          </h2>
          <iframe v-if="visibleIframe" src="./iframe.html"> </iframe>
          <button @click="iframePost">iframeのformからPOST</button>

          <ul class="memo">
            <li>（Chrome、Safari）cookie送信されない</li>
            <li>（FireFox）cookie送信される</li>
            <li>cookieにSameSite:Laxなし → cookie送信される</li>
            <li>cookieにSameSite:Lax → cookie送信されない</li>
          </ul>
        </section>
        <!-- formからGET -->
        <section>
          <h2>
            <span class="ok">OK</span><span>{{ urlUsers }}</span> へForm Get
          </h2>
          <form :action="urlUsers" method="GET">
            <input name="myvalue" value="isget" />
            <button type="submit">formからGET</button>
          </form>
          <ul class="memo">
            <li>cookieにSameSite:Laxなし → cookie送信される</li>
            <li>cookieにSameSite:Lax → cookie送信される</li>
            <li>ページ遷移と同様だから当然か</li>
            <li>Chrom、FireFox、Safariともに同じ</li>
          </ul>
        </section>
        <!-- section -->
        <section>
          <h2><span class="ok">OK</span>get users by fetch</h2>
          <button @click="getUsers">fetch getUsers</button>
          <ul>
            <li v-for="i in items" :key="`${i.name}`">
              <span>{{ i.name }}</span>
              <span>{{ i.email }}</span>
            </li>
          </ul>
          <ul class="memo">
            <li>Chrome、FireFox、Safariいずれも送信されない</li>
          </ul>
        </section>
        <!-- section -->
        <section>
          <h2>
            <span class="ok">OK</span>
            jsから
            <span>{{ urlUsers }}</span
            >へPOST
          </h2>
          <p>（X-XSRF-TOKENヘッダ付き）</p>
          <button @click="() => doPost(urlUsers)">doPostUsers</button>
          <span v-if="postResult[urlUsers]" class="postResult">{{ postResult[urlUsers] }}</span>
          <ul class="memo">
            <li>X-XSRF-TOKENヘッダ付きなのでプリフライトあり</li>
            <li>'Content-Type': 'application/json'なのでプリフライトあり</li>
            <li>Chrome、FireFox、Safariいずれも送信されない</li>
          </ul>
        </section>
        <!-- section -->
        <section v-if="false">
          <h2>
            jsから <span>{{ urlTop }}</span
            >へPOST
          </h2>
          <button @click="() => doPost(urlTop)">doPostTop</button>
          <span v-if="postResult[urlTop]" class="postResult">{{ postResult[urlTop] }}</span>
        </section>
        <!-- section -->
        <section>
          <h2>
            imgから <span>{{ urlTop }}</span
            >へGET
          </h2>
          <button @click="() => doGet()">doGet</button>
          <span>{{ visibleImg }}</span>
          <img v-if="visibleImg" :src="visibleImg" alt="" />

          <ul class="memo">
            <li>（Chrome、Safari）cookie送信されない</li>
            <li>（FireFox）cookie送信される</li>
            <li>cookieにSameSite:Laxなし → cookie送信される</li>
            <li>cookieにSameSite:Lax → cookie送信されない</li>
          </ul>
        </section>

        <h2>
          cookies
          <button @click="clearCookie">クリア</button>
          <p>
            <button @click="() => setNyao(true)">セット nyao with SameSite:Lax</button>
            <button @click="() => setNyao(false)">セット nyao with SameSiteなし</button>
          </p>
        </h2>
        <table>
          <thead>
            <th>name</th>
            <th>value</th>
          </thead>
          <tbody>
            <tr v-for="c in cookiesList" :key="c.name">
              <td>{{ c.name }}</td>
              <td>{{ c.value }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <article class="result">
      <div class="container">
        <h1>結論 以下が有効</h1>
        <h2>フォームからのPOST対策</h2>
        <ul>
          <li>フォームからpostするとドメイン関係なくcookie送信されるブラウザもあるが、CookieのSameSite:Laxというオプションをセットするとformでもcookieを送信しなくなる</li>
        </ul>
        <h2>プリフライトでのチェック</h2>
        <ul>
          <li>オリジナルヘッダをバックエンドで確認するようにする（中身はなんでもよい。存在が確認できればよい）</li>
          <li>オリジナルヘッダを入れればプリフライトが発生するため、攻撃サイトは必ずプリフライトを発生させることになる。</li>
          <li>プリフライトとは目的とするPOSTメソッドの直前に発生するOPTIONメソッドのことである（結果リクエストが2回走ることになるので少しレスポンスが遅くなる）</li>
          <li>プリフライト時のhttpメソッドはOPTIONなので、このリクエスト時にオリジンをチェックすれば許可していないサイトからのリクエストに対して対応できる</li>
        </ul>
        <h2>Cookie SameSite属性</h2>
        <ul>
          <li><b>（SameSite=Lax）</b>外部サイトからのGETメソッドでの遷移時(Top Level Navigation)のみ、Cookieを送信する。 これにより、Webサイトの使いやすさとセキュリティを両立させることができる</li>
          <li><b>（SameSite=Strict）</b>外部サイトからは常にCookieを送信しないためログイン済みのサイトに対して、外部サイトから遷移し直すと、未ログイン状態として処理される</li>
        </ul>
        <h2>Read CookieとWrite Cookieを用意するべき</h2>
        <ul>
          <li>
            <b>Read Cookie（SameSite=Lax）</b>
            ユーザのセッション維持するための Cookie 。これが送られればユーザがログイン状態と見なされる。しかし、それ以上の操作についてはこの Cookie では許可されない。
          </li>
          <li>
            <b>Write Cookie （SameSite=Strict）</b>
            書き込みを許可する Cookie 。例えばパスワードの変更や、投稿、購入、キャンセル、送金など、副作用を伴う操作の許可に使われる。
          </li>
        </ul>
      </div>
    </article>
  </main>
</template>
<!------------------------------->

<!------------------------------->
<script lang="ts">
import Vue from 'vue';
import Cookies from 'js-cookie';

// ----------------------
// default
// ----------------------
export default Vue.extend({
  data() {
    return {
      items: [],
      postResult: {},
      cookies: '',
      visibleImg: '',
      visibleIframe: false,
    };
  },
  computed: {
    urlUsers() {
      return 'http://localhost:8888/api/users/';
    },
    urlTop() {
      return 'http://localhost:8888/';
    },
    cookiesList() {
      const ret = [];
      for (const [name, value] of Object.entries(this.cookies)) {
        ret.push({ name, value });
      }
      console.log('ret', ret);
      return ret;
    },
  },
  created() {},
  mounted() {
    this.getUsers();
    this.updateCookie();
  },
  methods: {
    doGet() {
      this.visibleImg = `${this.urlTop}cancrow.png?rnd=${Math.random()}`;
      setTimeout(() => {
        this.visibleImg = '';
      }, 1000);
    },
    getUsers() {
      fetch(this.urlUsers)
        .then((response) => response.json())
        .then((data) => {
          this.items = data.users;
        });
    },
    async doPost(myurl) {
      console.log('doPost', this);
      const res = await fetch(myurl, {
        method: 'POST', // *GET, POST, PUT, DELETE, etc.
        mode: 'cors', // no-cors, *cors, same-origin
        cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
        credentials: 'include', // include, *same-origin, omit
        headers: {
          // 以下のどちらかでプリフライトは発生する
          'Content-Type': 'application/json',
          'X-XSRF-TOKEN': 'SECRET',
        },
        redirect: 'follow', // manual, *follow, error
        referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
        body: '', // body data type must match "Content-Type" header
      })
        .then((response) => response.json())
        .catch((response) => response.json());
      console.log('response', res);
      this.$set(this.postResult, myurl, JSON.stringify(res));
    },
    clearCookie() {
      const items = Cookies.get();
      for (const [name] of Object.entries(items)) {
        Cookies.remove(name, { path: '' }); // removed!
      }
      this.updateCookie();
    },
    setNyao(withLax) {
      if (withLax) {
        Cookies.set('nyao', `from-${location.hostname}-by-js`, { SameSite: 'Lax' });
      } else {
        Cookies.set('nyao', `from-${location.hostname}-by-js`);
      }

      this.updateCookie();
    },
    updateCookie() {
      this.cookies = Cookies.get();
    },
    iframePost() {
      this.visibleIframe = !this.visibleIframe;
      // this.$nextTick(() => {
      //   const $ref = this.$refs.iframeform;
      //   if($ref){
      //     $ref.submit();
      //   }
      // })
    },
  },
});
</script>
<!------------------------------->

<!------------------------------->
<style>
html {
  font-family: 'Helvetica Neue', Arial, 'Hiragino Kaku Gothic ProN', 'Hiragino Sans', Meiryo, sans-serif;
}
body {
  margin: 0;
}

header {
  text-align: center;
  background-color: #eee;
  padding: 3px 0;
  margin: 0;
}
header h1 {
  font-size: 16px;
  color: #999;
  font-weight: normal;
}
main {
  padding: 0 0 40px;
  font-size: 18px;
}
.container {
  width: 800px;
  margin: 0 auto;
}
section {
  border-bottom: dashed 1px #aaa;
  padding: 10px 0 10px;
}
h2 {
  font-size: 24px;
}
h2 span {
  margin: 0 10px;
  color: brown;
  vertical-align: middle;
}
input {
  font-size: 16px;
  padding: 10px 20px;
}
input:focus {
  outline: none;
}
button {
  font-size: 16px;
  margin-left: 20px;
  padding: 10px 20px;
  border-radius: 6px;
  box-shadow: 0 2px 2px 1px rgba(0, 0, 0, 0.2);
  border: none;
  color: black;
  cursor: pointer;
}
button:hover {
  background-color: #e2e2e2;
}
li span {
  margin-right: 1em;
}

table {
  width: 100%;
  border-collapse: collapse;
  background-color: #eee;
}
tr {
  border-top: solid 1px #fff;
}
th:nth-child(1) {
  border-right: solid 1px #fff;
}
td {
  padding: 10px;
  word-break: break-all;
}
td:nth-child(1) {
  width: 50%;
  border-right: solid 1px #fff;
}
.postResult {
  margin-left: 1em;
}
.memo {
  background-color: #eee;
  padding: 20px;
  font-size: 16px;
  border-radius: 4px;
  list-style: none;
}
b {
  color: darkred;
}
.warn {
  background-color: red;
  color: white;
  padding: 10px 10px;
  border-radius: 3px;
  font-size: 12px;
  margin-left: 1em;
}
.ok {
  background-color: green;
  color: white;
  padding: 10px 20px;
  border-radius: 3px;
  font-size: 12px;
  margin-left: 1em;
}
.try {
  margin: 0 0 20px;
}
.result {
  background-color: #e4e4e4;
  padding: 30px 0;
}
</style>