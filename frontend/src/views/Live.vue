<template>
  <div class="live">
    <h1 class="title">
      <router-link
              :to="{name: 'live', params: { channel: currentChannel.name }}"
      >{{ currentChannel.name }}</router-link>
    </h1>
    <div class="main">
      <ul class="channels">
        <router-link
                :class="{channel: true, active: channel.id === currentChannel.id}"
                :to="{name: 'live', params: { channel: channel.name }}"
                v-for="channel in sortedChannels" :key="channel.id">
        <li>
            #{{ channel.name }}
        </li>
        </router-link>
        <li class="padder"></li>
      </ul>
      <div class="messages" ref="messages" v-on:scroll="scrollHistoryUpdate">
        <div class="loadmore" v-if="channelMore(currentChannel.id)" v-on:click="loadMessagesBefore">
          Load more
        </div>
        <Message v-for="message in currentChannelMessages"
                 :message="message"
                 :key="message.id"
                 :channel-name="currentChannel.name"
                 :ref="'message' + message.id"
                 @loaded="scrollToBottom"
                 @scrollToPost="scrollToPost"
                 @insertNo="insertNo"
        ></Message>
      </div>
      <div v-if="scrollHistory" v-on:click="scrollToBottom(true)" class="scrollBackDown"></div>
    </div>
    <div class="inputs">
      <div class="row">
        <input placeholder="Username" aria-label="username" class="username text" type="text" v-model="userName"/>
        <textarea ref="bodyText" v-model="bodyText" v-on:keypress="keyhandle" :placeholder="'Message #' + currentChannel.name" aria-label="post" class="body text" rows="1"/>
        <button v-on:click="send" class="send text">Send</button>
      </div>
      <div class="row">
        <input v-model="userTrip" placeholder="Tripcode" aria-label="trip" class="trip text"/>
        <input type="file" class="file" ref="fileInput"/>
        <button v-on:click="clearFileInput" class="clear text">Remove chosen file</button>
      </div>
    </div>
  </div>
</template>

<script>
  import Vue from 'vue'
  import {mapActions, mapGetters} from "vuex";
  import Message from "../components/Message";

export default {
  name: 'live',
  components: {
    Message
  },
  data: () => ({
    scoped: {},
    scroll: false,
    scrollHistory: false,
    scrollBefore: false,
    scrollId: null,
  }),
  updated() {
    if (this.scrollId == null) {
      return;
    }
    if (this.currentChannelMessages.length === 0) {
      return;
    }
    const first = this.currentChannelMessages[0].id;
    if (first !== this.scrollId) {
      this.scrollToPost(this.scrollId);
      this.scrollId = null;
    }
  },
  watch: {
    currentChannelMessages() {
      this.scroll = !this.scrollHistory;
      if (this.scroll) {
        this.scrollToBottom();
      }
    },
    currentChannel() {
      this.scroll = true;
      this.scrollToBottom();
    },
    scrollBefore(val) {
      if (val) {
        if (this.currentChannelMessages.length === 0) {
          return;
        }
        this.loadMessagesBefore()
      }
    },
  },
  computed: {
    ...mapGetters('channels', [
      'sortedChannels',
      'currentChannel',
    ]),
    ...mapGetters('messages', [
      'messagesByChannel',
      'channelMore',
    ]),
    ...mapGetters('users', [
      'selfUsername',
      'selfTripcode',
    ]),
    currentChannelMessages() {
      return this.messagesByChannel(this.currentChannel.id);
    },
    userName: {
      get() {
        return this.selfUsername;
      },
      set(username) {
        this.$store.commit('users/setSelfUsername', {username});
      }
    },
    userTrip: {
      get() {
        return this.selfTripcode;
      },
      set(tripcode) {
        this.$store.commit('users/setSelfTripcode', {tripcode});
      }
    },
    bodyText: {
      get() {
        const ch = this.scoped[this.currentChannel.id];
        if (!ch) {
          return '';
        }
        return ch.body;
      },
      set(body) {
        const ch = this.scoped[this.currentChannel.id] || {};
        ch.body = body;
        Vue.set(this.scoped, this.currentChannel.id, ch);
      }
    }
  },
  methods: {
    ...mapActions([
      'sendPost',
      'sendCurrentChannel',
      'loadBefore',
    ]),
    scrollToBottom(force) {
      if (this.scroll || force) {
        this.$nextTick(function () {
          const el = this.$refs.messages;
          el.scroll(0, el.scrollHeight);
        });
      }
    },
    scrollHistoryUpdate() {
      const el = this.$refs.messages;
      let max = el.scrollHeight - el.offsetHeight;
      if (el.children.length > 0) {
        max -= el.children[el.children.length - 1].clientHeight;
      }
      this.scrollHistory = el.scrollTop < max;
      this.scrollBefore = el.scrollTop < el.clientHeight * 2;
    },
    clearFileInput() {
      this.$refs.fileInput.value = '';
    },
    clearPost() {
      Vue.set(this.scoped, this.currentChannel.id, {});
      this.clearFileInput();
    },
    keyhandle(key) {
      if (key.code === 'Enter' && !key.shiftKey) {
        key.preventDefault();
        this.send();
      }
    },
    send() {
      const ch = this.scoped[this.currentChannel.id] || {};
      this.sendPost({
        file: this.$refs.fileInput.files.length > 0 ? this.$refs.fileInput.files[0] : null,
        body: ch.body || null,
        trip: this.userTrip || null,
        name: this.userName || null,
        channel_id: this.currentChannel.id || null,
      }).then(() => this.clearPost()).catch((x) => {
        console.log(x.responseText);
      });
    },
    loadMessagesBefore() {
      if (this.channelMore(this.currentChannel.id)) {
        this.scrollId = this.currentChannelMessages[0].id;
        this.loadBefore({id: this.currentChannelMessages[0].id, channel_id: this.currentChannel.id});
      }
    },
    insertNo(id) {
      const el = this.$refs.bodyText;
      const ch = this.scoped[this.currentChannel.id] || {};
      ch.body = el.value.substring(0, el.selectionStart) + '>>' + id + ' ' + el.value.substring(el.selectionStart);
      el.value = ch.body;
    },
    scrollToPost(id) {
      console.log(id);
      const ref = this.$refs['message' + id];
      if (!ref || ref.length === 0) {
        return;
      }
      const el = ref[0].$el;
      this.$refs.messages.scroll(0, el.offsetTop);
    }
  }
}
</script>
<style>
  .live {
    display: flex;
    flex-direction: column;
    flex-grow: 1;
    overflow-y: hidden;
  }
  .live > .title {
    text-transform: capitalize;
    text-decoration: none;
    letter-spacing: -0.05ex;
  }
  .live > .title, .live > .inputs {
    flex-shrink: 0;
  }
  .live > .main {
    display: flex;
    flex-direction: row;
    flex-grow: 1;
    overflow-y: hidden;
    position: relative;
  }
  .live > .main > .channels {
    margin: 0;
    min-width: 10em;
    text-align: left;
    list-style: none;
    padding: 0;
    font-size: large;
    display: flex;
    flex-direction: column;
  }
  .live > .main > .channels > .channel {
    background: var(--bg-panel-color);
    padding: 2ex 1ex;
  }
  .live > .main > .channels > .active {
    background: var(--bg-panel-hover);
  }
  .live > .main > .channels > .channel:hover {
    background: var(--bg-panel-hover);
   }
  .live > .main > .channels > .padder {
    background: var(--bg-panel-color);
    flex-grow: 1;
    content: "";
    display: block;
  }
  .live > .main > .messages {
    background: var(--bg-panel-light);
    overflow-y: scroll;
    display: flex;
    flex-direction: column;
    justify-items: flex-start;
    flex-grow: 1;
  }

  .live > .main > .messages > .loadmore {
    color: var(--fg-link-color);
    cursor: pointer;
    padding: 1em;
  }

  .live > .inputs {
    border-radius: 0 0 1ex 1ex;
    padding: 1em;
    margin: 1em 0;
    background: var(--bg-form-color);
    display: flex;
    flex-direction: column
  }
  .live > .inputs > .row > *:first-child {
    width: 20%;
    border-left: 0.2ex solid var(--bg-accent-color);
  }
  .live > .inputs > .row > * {
    border: none;
    padding: 1ex;
    border-right: 0.2ex solid var(--bg-accent-color);
  }
  .live > .inputs > .row:first-child {
    border-top: 0.2ex solid var(--bg-accent-color);
  }
  .live > .inputs > .row {
    flex: 1;
    display: flex;
    flex-direction: row;
    border-bottom: 0.2ex solid var(--bg-accent-color);
  }
  .row:first-of-type {
    flex-wrap: wrap;
  }
  .live > .inputs > .row > .text {
    background: transparent;
    border-top: none;
    border-bottom: none;
    font: inherit;
    color: var(--fg-text-color);
  }
  .live > .inputs > .row > .username {
    text-align: center;
  }
  .live > .inputs > .row > .trip {
    text-align: center;
  }
  .live > .inputs > .row > .body {
    flex-grow: 1;
    resize: none;
  }
  .live > .inputs > .row >.file {
    flex-grow: 1;
  }

  .live > .inputs > .row > button {
    cursor: pointer;
    padding: 0 2em;
  }
  .live > .inputs > .row > button:hover {
    background: var(--bg-panel-hover);
  }
  .scrollBackDown {
    position: absolute;
    width: 2em;
    height: 2em;
    cursor: pointer;
    right: 2em;
    bottom: 0.5em;
  }
  .scrollBackDown:after,
  .scrollBackDown:before {
    content: '';
    display: block;
    position: absolute;
    height: 0.2em;
    width: 2em;
    background: var(--fg-link-color);
  }
  .scrollBackDown:after {
    top: 1em;
    left: 0.6em;
    transform: rotate(-45deg);
  }
  .scrollBackDown:before {
    top: 1em;
    left: -0.7em;
    transform: rotate(45deg);
  }
</style>