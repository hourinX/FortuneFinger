<template>
  <main class="login-page">
    <section class="visual-panel">
      <div class="visual-copy">
        <div class="english-title">F o r t u n e F i n g e r</div>
        <div class="english-subtitle">A n n u a l&nbsp;&nbsp;L u c k y&nbsp;&nbsp;D r a w</div>
        <h1>
          让每一次相聚<br /><strong>都充满<span>惊喜</span></strong>
        </h1>
        <i class="short-line" />
        <p>简单 · 有趣 · 高效<br /><small>一站式年会抽奖解决方案</small></p>
      </div>
      <div class="visual-footer">GOOD PEOPLE<br />BRING GOOD LUCK</div>
    </section>

    <section class="login-side">
      <div class="login-card">
        <div class="journey">开启精彩的抽奖之旅</div>
        <div class="brand">
          <img class="brand-icon" :src="brandIcon" alt="" />
          <img class="brand-name" :src="brandWordmark" alt="FortuneFinger" />
          <div class="brand-subtitle">后台管理系统</div>
          <i />
        </div>
        <LoginForm />
        <div class="values"><span />团结 · 感恩 · 展望 · 共赢<span /></div>
        <footer>© {{ year }} FortuneFinger. 让每一次年会，都值得期待。</footer>
      </div>
    </section>
  </main>
</template>

<script setup>
import brandIcon from '@/assets/branding/login-fortune-mark.png'
import brandWordmark from '@/assets/branding/fortune-finger-wordmark.png'
import LoginForm from './components/LoginForm.vue'

const year = new Date().getFullYear()
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 52.4% 47.6%;
  overflow: hidden;
  background: #f4f3f2;
}
.visual-panel {
  position: relative;
  min-height: 100vh;
  overflow: hidden;
  color: #fff;
  background-color: #681713;
  background-image: url('@/assets/images/login-side-bg.png');
  background-position: center;
  background-size: cover;
}
.visual-panel::after {
  position: absolute;
  inset: 0;
  content: '';
  background: linear-gradient(90deg, rgba(25, 4, 4, 0.32), rgba(69, 8, 8, 0.05));
}
.visual-copy {
  position: absolute;
  z-index: 2;
  top: 8.8%;
  left: 11%;
  padding-left: 35px;
  border-left: 1px solid rgba(255, 255, 255, 0.35);
}
.english-title {
  font-size: 17px;
  letter-spacing: 0.28em;
}
.english-subtitle {
  margin-top: 13px;
  color: rgba(255, 255, 255, 0.62);
  font-size: 12px;
  letter-spacing: 0.28em;
}
.visual-copy h1 {
  margin: 76px 0 0 12px;
  font-size: clamp(31px, 3vw, 48px);
  font-weight: 400;
  line-height: 1.48;
  letter-spacing: 0.01em;
}
.visual-copy h1 strong {
  font-size: 1.12em;
  font-weight: 700;
}
.visual-copy h1 span {
  color: var(--color-primary-hover);
}
.short-line {
  display: block;
  width: 36px;
  height: 4px;
  margin: 50px 0 24px 12px;
  background: var(--color-primary-border);
  border-radius: 4px;
}
.visual-copy p {
  margin-left: 12px;
  font-size: 17px;
  line-height: 1.7;
}
.visual-copy small {
  font-size: 14px;
}
.visual-footer {
  position: absolute;
  z-index: 2;
  bottom: 7%;
  left: 11%;
  padding-left: 35px;
  color: rgba(255, 255, 255, 0.72);
  font-size: 10px;
  line-height: 2;
  letter-spacing: 0.42em;
  border-left: 1px solid rgba(255, 255, 255, 0.35);
}
.login-side {
  min-height: 100vh;
  padding: 24px 26px 24px 0;
  background: linear-gradient(135deg, #f4f3f2, #ecebea);
}
.login-card {
  position: relative;
  min-height: calc(100vh - 48px);
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 42px 13% 48px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 0 22px 22px 0;
  box-shadow: 16px 0 45px rgba(40, 35, 34, 0.08);
}
.journey {
  width: 100%;
  color: #a3a3a3;
  font-size: 14px;
  text-align: right;
}
.brand {
  margin-top: 54px;
  text-align: center;
}
.brand-icon {
  display: block;
  width: 88px;
  height: 88px;
  margin: 0 auto;
  object-fit: contain;
}
.brand-name {
  display: block;
  width: clamp(230px, 19vw, 310px);
  height: auto;
  margin-top: 7px;
  margin-right: auto;
  margin-left: auto;
  object-fit: contain;
}
.brand-subtitle {
  margin-top: 9px;
  color: #686868;
  font-size: 20px;
  letter-spacing: 0.06em;
}
.brand i {
  display: block;
  width: 30px;
  height: 2px;
  margin: 23px auto 42px;
  background: var(--color-primary);
}
.values {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15px;
  margin-top: auto;
  color: #a2a2a2;
  font-size: 13px;
  white-space: nowrap;
}
.values span {
  width: 60px;
  height: 1px;
  background: #d8d8d8;
}
footer {
  margin-top: 53px;
  color: #aaa;
  font-size: 12px;
  text-align: center;
}
@media (max-width: 880px) {
  .login-page {
    display: block;
    background: linear-gradient(145deg, var(--color-sider-deep), var(--color-sider));
  }
  .visual-panel {
    display: none;
  }
  .login-side {
    min-height: 100vh;
    padding: 18px;
    background: transparent;
  }
  .login-card {
    min-height: calc(100vh - 36px);
    padding: 28px clamp(24px, 8vw, 64px);
    border-radius: 20px;
  }
  .journey {
    font-size: 12px;
  }
  .brand {
    margin-top: 35px;
  }
}
@media (min-width: 881px) and (max-height: 820px) {
  .login-card {
    padding-top: 28px;
    padding-bottom: 28px;
  }
  .brand {
    margin-top: 18px;
  }
  .brand-icon {
    width: 70px;
    height: 70px;
  }
  .brand i {
    margin: 15px auto 25px;
  }
  footer {
    margin-top: 24px;
  }
}
</style>
