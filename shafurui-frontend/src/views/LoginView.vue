<template>
  <main class="login-page">
    <section class="brand-panel" aria-label="视频相册">
      <div class="brand-content">
        <div class="brand">
          <div class="mark" aria-hidden="true">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
              <path
                d="M6 5.5h12a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2v-9a2 2 0 0 1 2-2Z"
                stroke="currentColor"
                stroke-width="1.8"
              />
              <path d="m10.2 9 4.7 3-4.7 3V9Z" fill="currentColor" />
            </svg>
          </div>
          <div>
            <h1>视频相册</h1>
            <p>Private Video Archive</p>
          </div>
        </div>

        <div class="hero-copy">
          <h2>进入你的私有影像库</h2>
          <p>用于家庭视频、手机导出素材和服务器静态目录的轻量访问入口。</p>
        </div>

        <div class="system-strip" aria-label="系统状态">
          <div class="metric">
            <b>/videos/</b>
            <span>静态文件前缀</span>
          </div>
          <div class="metric">
            <b>Range</b>
            <span>浏览器原生播放</span>
          </div>
          <div class="metric">
            <b>Basic</b>
            <span>建议叠加 nginx 认证</span>
          </div>
        </div>
      </div>
    </section>

    <section class="login-panel" aria-label="登录">
      <div class="login-card">
        <div class="card-head">
          <div class="status"><span class="status-dot"></span> LOCAL SESSION</div>
          <h2>登录</h2>
          <p>验证通过后进入按日期分组的视频相册。</p>
        </div>

        <form class="login-form" @submit.prevent="submit">
          <div class="field">
            <label for="account">账号</label>
            <div class="control">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                <path
                  d="M20 21a8 8 0 0 0-16 0M12 12a4 4 0 1 0 0-8 4 4 0 0 0 0 8Z"
                  stroke="currentColor"
                  stroke-width="1.9"
                  stroke-linecap="round"
                />
              </svg>
              <input v-model="account" id="account" type="text" autocomplete="username" />
            </div>
          </div>

          <div class="field">
            <label for="password">密码</label>
            <div class="control">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                <path
                  d="M7 10V8a5 5 0 0 1 10 0v2M6 10h12a1 1 0 0 1 1 1v8a1 1 0 0 1-1 1H6a1 1 0 0 1-1-1v-8a1 1 0 0 1 1-1Z"
                  stroke="currentColor"
                  stroke-width="1.9"
                  stroke-linecap="round"
                />
              </svg>
              <input
                v-model="password"
                id="password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
              />
              <button
                class="toggle-pass"
                type="button"
                title="显示或隐藏密码"
                aria-label="显示或隐藏密码"
                @click="showPassword = !showPassword"
              >
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                  <path
                    d="M2.5 12s3.5-6 9.5-6 9.5 6 9.5 6-3.5 6-9.5 6-9.5-6-9.5-6Z"
                    stroke="currentColor"
                    stroke-width="1.8"
                  />
                  <path
                    d="M12 15a3 3 0 1 0 0-6 3 3 0 0 0 0 6Z"
                    stroke="currentColor"
                    stroke-width="1.8"
                  />
                </svg>
              </button>
            </div>
          </div>

          <div class="row">
            <label class="check">
              <input v-model="remember" type="checkbox" />
              <span>保持登录</span>
            </label>
            <button class="link-button" type="button" @click="reset">重置输入</button>
          </div>

          <button class="primary" type="submit" :disabled="submitting">
            <span>{{ submitting ? "正在验证" : "登录" }}</span>
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path
                d="M5 12h14m-6-6 6 6-6 6"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </button>

          <div v-if="notice" class="notice" :class="noticeType">{{ notice }}</div>
        </form>
      </div>
    </section>
  </main>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { login } from "@/api/global";
import { useUserStore } from "@/stores/user";

defineOptions({
  name: "LoginView",
});

const router = useRouter();
const userStore = useUserStore();
const account = ref("");
const password = ref("");
const remember = ref(true);
const showPassword = ref(false);
const submitting = ref(false);
const notice = ref("");
const noticeType = ref("");

function reset() {
  account.value = "";
  password.value = "";
  noticeType.value = "";
  notice.value = "";
}

async function submit() {
  if (submitting.value) return;

  const username = account.value.trim();
  if (!username || !password.value) {
    noticeType.value = "error";
    notice.value = "请输入账号和密码";
    return;
  }

  submitting.value = true;
  noticeType.value = "";
  notice.value = "正在验证...";

  try {
    const tokens = await login({
      username,
      password: password.value,
    });
    userStore.setTokens(tokens);
    void userStore.fetchUserInfo(true).catch((error) => {
      console.warn("Failed to load user info after login.", error);
    });
    if (remember.value) {
      localStorage.setItem("albumSession", "remote");
    } else {
      localStorage.removeItem("albumSession");
    }
    noticeType.value = "success";
    notice.value = "验证通过，正在进入相册";
    await router.push("/");
  } catch (error) {
    console.warn("Failed to login.", error);
    noticeType.value = "error";
    notice.value = "账号或密码不正确";
  } finally {
    submitting.value = false;
  }
}
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Fira+Code:wght@500;600;700&family=Fira+Sans:wght@400;500;600;700&display=swap");

.login-page {
  --bg: #f6f5f2;
  --panel: #ffffff;
  --ink: #101013;
  --muted: #68665f;
  --line: #d9d6cf;
  --line-strong: #b9b5aa;
  --dark: #171719;
  --accent: #c43f2f;
  --teal: #1f7a78;
  --good: #18794e;
  --shadow: 0 22px 70px rgba(18, 18, 20, 0.16);
  --radius: 8px;
  min-height: 100vh;
  display: grid;
  grid-template-columns: minmax(420px, 0.92fr) minmax(420px, 1.08fr);
  color: var(--ink);
  background:
    linear-gradient(90deg, rgba(16, 16, 19, 0.045) 1px, transparent 1px) 0 0 / 44px 44px,
    linear-gradient(rgba(16, 16, 19, 0.035) 1px, transparent 1px) 0 0 / 44px 44px,
    var(--bg);
  font-family: "Fira Sans", "PingFang SC", "Microsoft YaHei", sans-serif;
}

button,
input {
  font: inherit;
}

button {
  cursor: pointer;
}

.brand-panel {
  position: relative;
  min-height: 100vh;
  padding: 34px;
  color: #fff;
  background:
    radial-gradient(circle at 16% 18%, rgba(196, 63, 47, 0.38), transparent 28%),
    linear-gradient(135deg, #151518 0%, #25252a 54%, #111113 100%);
  overflow: hidden;
}

.brand-panel::before {
  content: "";
  position: absolute;
  inset: 0;
  background:
    repeating-linear-gradient(90deg, rgba(255, 255, 255, 0.14) 0 1px, transparent 1px 30px),
    repeating-linear-gradient(0deg, rgba(255, 255, 255, 0.08) 0 1px, transparent 1px 30px);
  mask-image: linear-gradient(120deg, rgba(0, 0, 0, 0.85), transparent 72%);
}

.brand-panel::after {
  content: "";
  position: absolute;
  right: -120px;
  bottom: -80px;
  width: 520px;
  height: 300px;
  background:
    linear-gradient(
      135deg,
      transparent 0 36%,
      rgba(255, 255, 255, 0.22) 37% 43%,
      transparent 44% 100%
    ),
    linear-gradient(90deg, rgba(31, 122, 120, 0.68), rgba(196, 63, 47, 0.38));
  transform: rotate(-8deg);
  opacity: 0.78;
}

.brand-content {
  position: relative;
  z-index: 1;
  min-height: calc(100vh - 68px);
  display: grid;
  grid-template-rows: auto 1fr auto;
  gap: 36px;
}

.brand,
.row,
.check {
  display: flex;
  align-items: center;
}

.brand {
  gap: 14px;
}

.mark {
  width: 48px;
  height: 48px;
  border-radius: var(--radius);
  display: grid;
  place-items: center;
  color: #fff;
  background:
    linear-gradient(135deg, transparent 0 42%, rgba(255, 255, 255, 0.28) 43% 56%, transparent 57%),
    rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.18);
}

.brand h1,
.hero-copy h2,
.card-head h2 {
  margin: 0;
  letter-spacing: 0;
}

.brand h1 {
  font-size: 22px;
  line-height: 1.1;
}

.brand p,
.hero-copy p {
  margin: 5px 0 0;
  color: rgba(255, 255, 255, 0.62);
}

.hero-copy {
  align-self: center;
  max-width: 580px;
}

.hero-copy h2 {
  font-size: clamp(42px, 6vw, 76px);
  line-height: 0.94;
}

.hero-copy p {
  margin-top: 22px;
  max-width: 500px;
  font-size: 16px;
  line-height: 1.7;
}

.system-strip {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
  max-width: 580px;
}

.metric {
  min-height: 92px;
  padding: 14px;
  border-radius: var(--radius);
  border: 1px solid rgba(255, 255, 255, 0.16);
  background: rgba(255, 255, 255, 0.07);
  backdrop-filter: blur(10px);
}

.metric b {
  display: block;
  font-family: "Fira Code", monospace;
  font-size: 21px;
}

.metric span {
  display: block;
  margin-top: 9px;
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
  line-height: 1.4;
}

.login-panel {
  min-height: 100vh;
  padding: 34px;
  display: grid;
  place-items: center;
}

.login-card {
  width: min(100%, 470px);
  border: 1px solid var(--line);
  border-radius: var(--radius);
  background: rgba(255, 255, 255, 0.82);
  box-shadow: var(--shadow);
  overflow: hidden;
}

.card-head {
  padding: 24px 24px 18px;
  border-bottom: 1px solid var(--line);
}

.status {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-height: 30px;
  padding: 0 10px;
  border-radius: 999px;
  color: var(--good);
  background: rgba(24, 121, 78, 0.1);
  border: 1px solid rgba(24, 121, 78, 0.22);
  font-size: 12px;
  font-weight: 700;
}

.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: currentColor;
}

.card-head h2 {
  margin-top: 18px;
  font-size: 30px;
  line-height: 1.1;
}

.card-head p {
  margin: 10px 0 0;
  color: var(--muted);
  line-height: 1.55;
}

.login-form {
  padding: 22px 24px 24px;
  display: grid;
  gap: 16px;
}

.field {
  display: grid;
  gap: 8px;
}

.field label {
  color: var(--dark);
  font-size: 13px;
  font-weight: 700;
}

.control {
  position: relative;
}

.control svg {
  position: absolute;
  left: 14px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--muted);
  pointer-events: none;
}

.control input {
  width: 100%;
  height: 48px;
  border-radius: var(--radius);
  border: 1px solid var(--line-strong);
  background: #fff;
  color: var(--ink);
  outline: none;
  padding: 0 48px 0 44px;
}

.control input:focus {
  border-color: var(--dark);
  box-shadow: 0 0 0 3px rgba(16, 16, 19, 0.12);
}

.toggle-pass {
  position: absolute;
  right: 6px;
  top: 50%;
  transform: translateY(-50%);
  width: 38px;
  height: 38px;
  border: 0;
  border-radius: 7px;
  color: var(--dark);
  background: transparent;
  display: grid;
  place-items: center;
}

.toggle-pass svg {
  position: static;
  transform: none;
  color: inherit;
}

.toggle-pass:hover {
  background: rgba(16, 16, 19, 0.07);
}

.row {
  justify-content: space-between;
  gap: 12px;
}

.check {
  gap: 8px;
  color: var(--muted);
  font-size: 13px;
}

.check input {
  width: 16px;
  height: 16px;
  accent-color: var(--dark);
}

.link-button {
  border: 0;
  padding: 0;
  background: transparent;
  color: var(--teal);
  font-size: 13px;
  font-weight: 700;
}

.primary {
  height: 50px;
  border-radius: var(--radius);
  border: 1px solid var(--dark);
  background: var(--dark);
  color: #fff;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  font-weight: 700;
  transition:
    transform 180ms ease,
    background 180ms ease,
    box-shadow 180ms ease;
}

.primary:hover {
  transform: translateY(-1px);
  background: #050506;
  box-shadow: 0 14px 30px rgba(18, 18, 20, 0.2);
}

.primary:disabled {
  cursor: wait;
  opacity: 0.74;
  transform: none;
}

.notice {
  min-height: 42px;
  border-radius: var(--radius);
  border: 1px solid var(--line);
  background: rgba(246, 245, 242, 0.78);
  color: var(--muted);
  padding: 11px 12px;
  font-size: 13px;
  line-height: 1.45;
}

.notice.error {
  color: var(--accent);
  background: rgba(196, 63, 47, 0.08);
  border-color: rgba(196, 63, 47, 0.28);
}

.notice.success {
  color: var(--good);
  background: rgba(24, 121, 78, 0.09);
  border-color: rgba(24, 121, 78, 0.26);
}

.mini-link {
  color: var(--dark);
  text-decoration: none;
  font-weight: 700;
}

@media (max-width: 980px) {
  .login-page {
    grid-template-columns: 1fr;
  }

  .brand-panel {
    min-height: auto;
    padding: 26px;
  }

  .brand-content {
    min-height: 460px;
  }

  .login-panel {
    min-height: auto;
    padding: 28px 20px 44px;
  }
}

@media (max-width: 620px) {
  .brand-panel {
    padding: 20px;
  }

  .brand-content {
    min-height: 420px;
  }

  .hero-copy h2 {
    font-size: 40px;
  }

  .system-strip {
    grid-template-columns: 1fr;
  }

  .card-head,
  .login-form,
  .card-foot {
    padding-left: 18px;
    padding-right: 18px;
  }

  .row,
  .card-foot {
    align-items: flex-start;
    display: flex;
    flex-direction: column;
  }
}
</style>
