<script setup lang="ts">
import { nextTick, ref } from 'vue'
import AppTabBar from '@/components/app-tabbar/index.vue'
import BrandHeader from '@/views/checkin/components/BrandHeader.vue'
import ChatMessage, { type ChatMessageItem } from './components/ChatMessage.vue'

const draft = ref('')
const messageList = ref<HTMLElement>()
const messages = ref<ChatMessageItem[]>([
  { id: 1, name: '林清越', content: '今晚不见不散，祝大家都能抽到心仪的奖品！', time: '18:36', avatar: '林', onScreen: true },
  { id: 2, name: '周景行', content: '很高兴再次和大家相聚，新的一年一起加油。', time: '18:37', avatar: '周', onScreen: true },
  { id: 3, name: '许南星', content: '现场布置得太漂亮了，已经开始期待抽奖环节了。', time: '18:38', avatar: '许' },
  { id: 4, name: '活动助手', content: '欢迎大家入场，发送的留言可选择同步至现场大屏。', time: '18:39', avatar: '福' },
])

async function send() {
  const content = draft.value.trim()
  if (!content) return
  messages.value.push({
    id: Date.now(),
    name: '我',
    content,
    time: new Intl.DateTimeFormat('zh-CN', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    }).format(new Date()),
    avatar: '我',
    own: true,
    onScreen: true,
  })
  draft.value = ''
  await nextTick()
  messageList.value?.scrollTo({ top: messageList.value.scrollHeight, behavior: 'smooth' })
}
</script>

<template>
  <main class="chat-page">
    <header class="chat-header">
      <BrandHeader company="FortuneFinger" caption="LIVE MOMENTS" />
    </header>

    <section class="chat-title">
      <h1>弹幕聊天</h1>
      <strong><i></i>128<small>人在线</small></strong>
    </section>

    <aside class="chat-notice">
      <span>公告</span>
      <p>签到正在进行中，第一轮抽奖将在 19:30 开始</p>
    </aside>

    <section ref="messageList" class="chat-list">
      <ChatMessage v-for="message in messages" :key="message.id" :message="message" />
    </section>

    <footer class="chat-composer">
      <form @submit.prevent="send">
        <div class="chat-composer__input">
          <input v-model="draft" maxlength="60" placeholder="发送一条现场弹幕…" />
          <small>{{ draft.length }}/60</small>
        </div>
        <button type="submit" :disabled="!draft.trim()">发送</button>
      </form>
    </footer>
    <AppTabBar active="chat" />
  </main>
</template>

<style scoped>
.chat-page { display: grid; width: min(100%, 440px); height: 100dvh; margin: 0 auto; box-sizing: border-box; grid-template-rows: auto auto auto minmax(0, 1fr) auto auto; overflow: hidden; background: linear-gradient(180deg, #f8f4ed, #f2ece4); }
.chat-header { padding: max(24px, env(safe-area-inset-top)) 18px 12px; }
.chat-title { display: flex; min-height: 55px; padding: 4px 20px 12px; align-items: center; justify-content: space-between; }
.chat-title h1 { margin: 0; color: #302c29; font-size: 23px; }
.chat-title strong { display: flex; align-items: center; gap: 5px; color: #d9413a; font-size: 18px; font-variant-numeric: tabular-nums; }
.chat-title strong i { width: 6px; height: 6px; background: #56a970; border-radius: 50%; box-shadow: 0 0 8px rgb(86 169 112 / 65%); }
.chat-title strong small { color: #9f968f; font-size: 9px; font-weight: 400; }
.chat-notice { display: flex; margin: 0 18px 12px; padding: 10px 12px; align-items: center; gap: 9px; color: #756c65; background: rgb(255 255 255 / 70%); border-radius: 10px; }
.chat-notice span { flex: 0 0 auto; padding: 3px 6px; color: #fff; font-size: 8px; background: #df433c; border-radius: 5px; }
.chat-notice p { overflow: hidden; margin: 0; font-size: 10px; text-overflow: ellipsis; white-space: nowrap; }
.chat-list { display: flex; min-height: 0; padding: 8px 18px 24px; flex-direction: column; gap: 16px; overflow-y: auto; overscroll-behavior: contain; scrollbar-width: none; }
.chat-list::-webkit-scrollbar { display: none; }
.chat-composer { padding: 12px 14px; background: rgb(255 255 255 / 94%); border-top: 1px solid #ece4da; box-shadow: 0 -12px 35px rgb(75 54 40 / 7%); }
.chat-composer form { display: grid; grid-template-columns: minmax(0, 1fr) 64px; align-items: center; gap: 10px; }
.chat-composer__input { position: relative; overflow: hidden; background: #f5f1eb; border-radius: 14px; box-shadow: inset 0 0 0 1px transparent; transition: background .2s, box-shadow .2s; }
.chat-composer__input:focus-within { background: #fffaf7; box-shadow: inset 0 0 0 1px rgb(221 68 61 / 28%); }
.chat-composer form input { width: 100%; height: 48px; padding: 0 48px 0 14px; box-sizing: border-box; color: #332f2c; font-size: 16px; outline: none; background: transparent; border: 0; }
.chat-composer form input::placeholder { color: #aaa19a; font-size: 13px; }
.chat-composer__input small { position: absolute; right: 11px; bottom: 8px; color: #b5ada6; font-size: 8px; font-variant-numeric: tabular-nums; }
.chat-composer button { height: 48px; color: #fff; font-size: 13px; font-weight: 650; background: linear-gradient(135deg, #ec4b44, #d73431); border: 0; border-radius: 14px; }
.chat-composer button:disabled { color: #aaa19a; background: #ece7e1; }
</style>
