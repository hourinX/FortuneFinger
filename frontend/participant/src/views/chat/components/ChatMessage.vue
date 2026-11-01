<script setup lang="ts">
export interface ChatMessageItem {
  id: number
  name: string
  content: string
  time: string
  avatar: string
  own?: boolean
  onScreen?: boolean
}

defineProps<{ message: ChatMessageItem }>()
</script>

<template>
  <article class="chat-message" :class="{ 'chat-message--own': message.own }">
    <div class="chat-message__avatar">{{ message.avatar }}</div>
    <div class="chat-message__body">
      <div class="chat-message__meta">
        <strong>{{ message.name }}</strong>
        <span v-if="message.onScreen">已上墙</span>
        <time>{{ message.time }}</time>
      </div>
      <p>{{ message.content }}</p>
    </div>
  </article>
</template>

<style scoped>
.chat-message { display: flex; align-items: flex-start; gap: 10px; }
.chat-message__avatar { display: grid; width: 36px; height: 36px; flex: 0 0 36px; place-items: center; color: #fff; font-size: 12px; font-weight: 650; background: linear-gradient(145deg, #d8a15f, #b7743a); border-radius: 50%; }
.chat-message__body { min-width: 0; max-width: calc(100% - 58px); }
.chat-message__meta { display: flex; margin: 0 4px 5px; align-items: center; gap: 7px; }
.chat-message__meta strong { color: #6c635d; font-size: 11px; }
.chat-message__meta span { padding: 2px 5px; color: #d9433c; font-size: 8px; background: #fff0ed; border-radius: 5px; }
.chat-message__meta time { color: #b1a8a0; font-size: 9px; }
.chat-message p { margin: 0; padding: 10px 13px; color: #3f3935; font-size: 14px; line-height: 20px; background: #fff; border-radius: 4px 14px 14px; box-shadow: 0 6px 20px rgb(77 57 42 / 7%); }
.chat-message--own { flex-direction: row-reverse; }
.chat-message--own .chat-message__avatar { background: linear-gradient(145deg, #ef5a52, #d53633); }
.chat-message--own .chat-message__body { display: flex; flex-direction: column; align-items: flex-end; }
.chat-message--own .chat-message__meta { flex-direction: row-reverse; }
.chat-message--own p { color: #fff; background: linear-gradient(135deg, #ed4a43, #d73532); border-radius: 14px 4px 14px 14px; }
</style>
