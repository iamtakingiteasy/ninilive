import Vue from 'vue'
import Vuex from 'vuex'
import _ from 'lodash';


Vue.use(Vuex);

import {channelsModule} from './channels';
import {usersModule} from "./users";
import {messagesModule} from "./messages";
import {endpointHTTP} from "../config";

const messageInit = () => () => {
  console.log('connected');
};

const messageHandler = (context) => (msg) => {
  const payload = JSON.parse(msg.data);
  console.log(payload.kind, payload);
  switch (payload.kind) {
    case 'init':
      context.commit('users/self', payload.data);
      break;
    case 'sessions':
      context.commit('users/set', payload.data);
      break;
    case 'channels':
      context.commit('channels/set', payload.data);
      break;
    case 'channelAdd':
      context.commit('channels/add', {channel: payload.data});
      break;
    case 'channelUpdate':
      context.commit('channels/add', {channel: payload.data});
      break;
    case 'channelRemove':
      context.commit('channels/del', {channel: payload.data});
      break;
    case 'online':
      context.commit('users/add', {session: payload.data});
      break;
    case 'offline':
      context.commit('users/del', {session: payload.data});
      break;
    case 'active':
      context.dispatch('channels/active', {session: payload.data});
      break;
    case 'messages':
      if (payload.id === 'before' || payload.id === 'page') {
        payload.data.more = payload.data.more || false;
      } else {
        payload.data.more = payload.data.more || null;
      }
      context.commit('messages/add', payload.data);
      break;
    case 'messagesRemove':
      context.commit('messages/del', payload.data);
      break;
  }
  //console.log(msg, context);
};

const messageError = (context) => (msg) => {
  console.error(msg, context);
};

const messageClose = (context) => (msg) => {
  console.log(msg);
  context.dispatch('reconnect');
};

const createWS = _.throttle((context) => {
  const ws = new WebSocket(context.state.remote);
  ws.onmessage = messageHandler(context);
  ws.onerror = messageError(context);
  ws.onopen = messageInit(context);
  ws.onclose = messageClose(context);
  context.commit('setWS', {ws});
}, 5000);

export default new Vuex.Store({
  state: {
    ws: null,
    remote: null,
  },
  mutations: {
    setRemote(state, {remote}) {
      state.remote = remote;
    },
    setWS(state, {ws}) {
      state.ws = ws;
    }
  },
  actions: {
    connect(context, {remote}) {
      context.commit('setRemote', {remote});
      context.dispatch('reconnect');
    },
    reconnect(context) {
      createWS(context);
    },
    fileUpload(context, {file}) {
      return new Promise(function (resolve, reject) {
        const req = new XMLHttpRequest();
        req.open("POST", endpointHTTP + '/ws/upload');
        req.onload = () => {
          if (req.status / 100 === 2) {
            resolve(req.response);
          } else {
            reject(req)
          }
        };
        req.onerror = () => reject(req);
        const formData = new FormData();
        formData.append("file", file);
        formData.append("session", context.state.users.session);
        req.send(formData);
      })
    },
    async sendPost(context, {file, body, trip, name, channel_id}) {
      console.log(context.state.users.session, file, body, trip, name, channel_id);
      if (file == null && (body == null || body.trim().length === 0)) {
        return Promise.resolve();
      }
      const post = {
        body,
        trip,
        name,
        channel_id
      };
      if (file != null) {
        const res = await context.dispatch('fileUpload', {file});
        post['file'] = {
          path: res,
          name: file.name,
        }
      }
      context.state.ws.send(JSON.stringify({
        id: '',
        kind: 'messageSend',
        data: post,
      }));
    },
    async sendCurrentChannel(context, {id}) {
      context.state.ws.send(JSON.stringify({
        id: '',
        kind: 'channelSelect',
        data: {id},
      }));
    },
    async loadBefore(context, {id, channel_id}) {
      context.state.ws.send(JSON.stringify({
        id: 'before',
        kind: 'messageBefore',
        data: {id, channel_id, limit: 100},
      }));
    }
  },
  modules: {
    channels: channelsModule,
    users: usersModule,
    messages: messagesModule,
  }
});