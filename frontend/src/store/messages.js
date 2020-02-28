import Vue from 'vue'
import Long from "long";

export const messagesModule = {
    namespaced: true,
    state: {
        channels: {}
    },
    mutations: {
        set(state, {channel_id, messages = [], more}) {
            const ch = state.channels[channel_id] || {
                messages: {},
                ids: [],
                more: more || false,
            };
            ch.messages = messages.reduce((acc, v) => ({...acc, [v.id]: v}), {});
            ch.ids = messages.map(x => Long.fromString(x.id, true)).sort((a, b) => a.compare(b));
            Vue.set(state.channels, channel_id, ch);
        },
        add(state, {channel_id, messages = [], more}) {
            const ch = state.channels[channel_id] || {
                messages: [],
                ids: [],
                more: more || false,
            };
            if (more === false) {
                ch.more = false;
            }
            ch.messages = messages.reduce((acc, v) => ({...acc, [v.id]: v}), ch.messages);
            const longs = messages.map(x => Long.fromString(x.id, true));
            const newids = longs.filter(x => !ch.ids.some(y => y.equals(x)));
            ch.ids = [...ch.ids, ...newids].sort((a, b) => a.compare(b));
            Vue.set(state.channels, channel_id, ch);
        },
        del(state, {channel_id, messages = []}) {
            const ch = state.channels[channel_id];
            if (!ch) {
                return
            }
            const longs = messages.map(x => Long.fromString(x.id, true));
            ch.ids = ch.ids.filter(x => !longs.some(y => y.equals(x)));
            Vue.set(state.channels, channel_id, ch);
        }
    },
    getters: {
        channelMore(state) {
            return function (id) {
                const ch = state.channels[id];
                if (!ch) {
                    return false;
                }
                return ch.more;
            }
        },
        messagesByChannel(state) {
            return function (id) {
                const ch = state.channels[id];
                if (!ch) {
                    return [];
                }
                return ch.ids.map(x => ch.messages[x.toString()]);
            }
        }
    }
};