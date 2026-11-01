<script setup lang="ts">
import AppTabBar from '@/components/app-tabbar/index.vue'
import BrandHeader from '@/views/checkin/components/BrandHeader.vue'

const participant = {
  name: '许南星',
  company: '苏州研发中心',
  phone: '138 **** 2026',
  checkinTime: '18:32',
}
</script>

<template>
  <main class="profile-page">
    <header class="profile-header">
      <BrandHeader company="FortuneFinger" caption="MY EVENT" />
    </header>

    <section class="profile-scroll">
      <section class="identity-card">
        <div class="identity-card__avatar">许</div>
        <div class="identity-card__content">
          <p>今晚，欢迎你</p>
          <h1>{{ participant.name }}</h1>
          <span>{{ participant.company }}</span>
        </div>
        <div class="identity-card__state">
          <i>✓</i>
          <span>已签到</span>
        </div>
      </section>

      <section class="event-card">
        <div class="event-card__topline">
          <span>当前活动</span>
          <strong>进行中</strong>
        </div>
        <h2>2026 FortuneFinger 年度盛典</h2>
        <p>让每一次相聚，都充满惊喜</p>
        <div class="event-card__meta">
          <span><small>签到时间</small>{{ participant.checkinTime }}</span>
          <span><small>活动日期</small>10 月 31 日</span>
        </div>
      </section>

      <section class="profile-section">
        <header>
          <div>
            <p>签到资料</p>
            <span>用于确认现场参与身份</span>
          </div>
          <i>已确认</i>
        </header>
        <dl>
          <div>
            <dt>真实姓名</dt>
            <dd>{{ participant.name }}</dd>
          </div>
          <div>
            <dt>所属企业</dt>
            <dd>{{ participant.company }}</dd>
          </div>
          <div>
            <dt>手机号码</dt>
            <dd>{{ participant.phone }}</dd>
          </div>
        </dl>
      </section>

      <section class="profile-section prize-section">
        <header>
          <div>
            <p>我的中奖记录</p>
            <span>中奖结果公布后将在这里同步</span>
          </div>
        </header>
        <div class="prize-section__empty">
          <div class="prize-section__gift"><i></i></div>
          <div>
            <strong>幸运正在路上</strong>
            <span>请留意现场抽奖结果</span>
          </div>
        </div>
      </section>
    </section>

    <AppTabBar active="profile" />
  </main>
</template>

<style scoped>
.profile-page { display: grid; width: min(100%, 440px); height: 100dvh; margin: 0 auto; box-sizing: border-box; grid-template-rows: auto minmax(0, 1fr) auto; overflow: hidden; background: linear-gradient(180deg, #f8f4ed, #f2ece4); }
.profile-header { padding: max(24px, env(safe-area-inset-top)) 18px 12px; }
.profile-scroll { min-height: 0; padding: 8px 18px 26px; overflow-y: auto; overscroll-behavior: contain; scrollbar-width: none; }
.profile-scroll::-webkit-scrollbar { display: none; }
.identity-card { position: relative; display: flex; min-height: 102px; padding: 22px; box-sizing: border-box; align-items: center; gap: 14px; overflow: hidden; background: linear-gradient(135deg, #42332d, #2c2725); border-radius: 22px; box-shadow: 0 20px 45px rgb(68 45 34 / 18%); }
.identity-card::after { position: absolute; right: -34px; bottom: -54px; width: 140px; height: 140px; background: radial-gradient(circle, rgb(229 91 76 / 26%), transparent 68%); content: ''; }
.identity-card__avatar { display: grid; width: 56px; height: 56px; flex: 0 0 56px; place-items: center; color: #fff; font-size: 19px; font-weight: 650; background: linear-gradient(145deg, #ef6259, #cf332f); border: 2px solid rgb(255 255 255 / 28%); border-radius: 50%; }
.identity-card__content { min-width: 0; }
.identity-card__content p { margin: 0 0 2px; color: #c5b9b2; font-size: 9px; }
.identity-card__content h1 { margin: 0; color: #fff; font-size: 22px; font-weight: 650; }
.identity-card__content span { display: block; margin-top: 4px; overflow: hidden; color: #b8ada7; font-size: 10px; text-overflow: ellipsis; white-space: nowrap; }
.identity-card__state { position: relative; z-index: 1; display: flex; margin-left: auto; flex: 0 0 auto; align-items: center; gap: 5px; color: #e9d5c4; font-size: 9px; }
.identity-card__state i { display: grid; width: 17px; height: 17px; place-items: center; color: #fff; font-size: 9px; font-style: normal; background: #58a36e; border-radius: 50%; }
.event-card { margin-top: 14px; padding: 19px 20px 17px; background: rgb(255 255 255 / 90%); border: 1px solid rgb(255 255 255 / 80%); border-radius: 18px; box-shadow: 0 14px 36px rgb(79 58 43 / 7%); }
.event-card__topline { display: flex; align-items: center; justify-content: space-between; }
.event-card__topline > span { color: #9b9189; font-size: 9px; }
.event-card__topline strong { color: #5c9a6d; font-size: 9px; font-weight: 600; }
.event-card h2 { margin: 11px 0 5px; color: #332e2b; font-size: 17px; }
.event-card > p { margin: 0; color: #91877f; font-size: 10px; }
.event-card__meta { display: grid; margin-top: 17px; padding-top: 14px; grid-template-columns: repeat(2, 1fr); border-top: 1px solid #eee8e1; }
.event-card__meta span { color: #4b4541; font-size: 12px; font-weight: 600; }
.event-card__meta span + span { padding-left: 18px; border-left: 1px solid #eee8e1; }
.event-card__meta small { display: block; margin-bottom: 4px; color: #aaa19a; font-size: 8px; font-weight: 400; }
.profile-section { margin-top: 14px; padding: 18px 20px; background: rgb(255 255 255 / 84%); border-radius: 18px; }
.profile-section > header { display: flex; align-items: flex-start; justify-content: space-between; }
.profile-section header p { margin: 0; color: #3f3935; font-size: 14px; font-weight: 650; }
.profile-section header span { display: block; margin-top: 4px; color: #aaa19a; font-size: 9px; }
.profile-section header > i { padding: 4px 7px; color: #568465; font-size: 8px; font-style: normal; background: #edf6ef; border-radius: 8px; }
.profile-section dl { margin: 14px 0 0; }
.profile-section dl div { display: flex; min-height: 41px; align-items: center; justify-content: space-between; border-top: 1px solid #f0ebe5; }
.profile-section dt { color: #9d938b; font-size: 10px; }
.profile-section dd { margin: 0; color: #4b4541; font-size: 11px; }
.prize-section__empty { display: flex; margin-top: 15px; padding-top: 15px; align-items: center; gap: 12px; border-top: 1px solid #f0ebe5; }
.prize-section__gift { position: relative; width: 38px; height: 34px; flex: 0 0 38px; background: #fff0ed; border-radius: 8px; }
.prize-section__gift::before { position: absolute; top: 11px; right: 5px; left: 5px; height: 16px; background: #e65049; border-radius: 2px 2px 5px 5px; content: ''; }
.prize-section__gift::after { position: absolute; top: 8px; left: 3px; width: 32px; height: 6px; background: #d83d37; border-radius: 2px; content: ''; }
.prize-section__gift i { position: absolute; z-index: 1; top: 8px; left: 17px; width: 4px; height: 19px; background: #ffd8b2; }
.prize-section__empty strong { display: block; color: #544c47; font-size: 11px; }
.prize-section__empty span { display: block; margin-top: 3px; color: #aaa19a; font-size: 9px; }
</style>
