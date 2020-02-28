import Vue from 'vue'

export const channelsModule = {
    namespaced: true,
    state: {
        channels: [],
        selected: null,
        selectedId: null,
        selectedName: null,
    },
    mutations: {
        set(state, {channels = []}) {
            state.channels = channels.sort((a,b) => (a.order || 0) - (b.order || 0));
            state.selected = null;
            if (state.selectedId != null) {
                const c = state.channels.find(x => x.name === state.selectedId);
                if (c != null) {
                    state.selectedName = c.name;
                    state.selected = c;
                }
            } else if (state.selectedName != null) {
                const c = state.channels.find(x => x.name === state.selectedName);
                if (c != null) {
                    state.selectedId = c.id;
                    state.selected = c;
                }
            }
        },
        add(state, {channel = {}}) {
            const idx = state.channels.findIndex(e => e.id === channel.id);
            if (idx !== -1) {
                Vue.set(state.channels, idx, channel);
            } else {
                state.channels.push(channel);
            }
            if (state.selected && state.selected.id === channel.id) {
                state.selected = channel;
            }
        },
        del(state, {channel = {}}) {
            const idx = state.channels.findIndex(e => e.id === channel.id);
            if (idx !== -1) {
                state.channels.splice(idx, 1);
            }
            if (state.selected && state.selected.id === channel.id) {
                state.selected = null;
            }
        },
    },
    actions: {
        selectNone(context) {
            context.state.selected = null;
        },
        selectId(context, {id}) {
            context.state.selectedId = id;
            const c = context.state.channels.find(x => x.name === id);
            if (c != null) {
                context.state.selectedName = c.name;
                context.state.selected = c;
                context.dispatch('sendCurrentChannel', {id: id}, {root: true});
            }
        },
        selectName(context, {name}) {
            context.state.selectedName = name;
            const c = context.state.channels.find(x => x.name === name);
            if (c != null) {
                context.state.selectedId = c.id;
                context.state.selected = c;
                context.dispatch('sendCurrentChannel', {id: c.id}, {root: true});
            }
        },
        active(context, {session: {active}}) {
            Object.entries(active).forEach(([k,v]) => {
                const idx = context.state.channels.findIndex(e => e.id === k);
                if (idx > -1) {
                    const channel = context.state.channels[idx];
                    channel.active = v;
                    Vue.set(context.state.channels, idx, channel);
                }
            })
        }
    },
    getters: {
        sortedChannels(state) {
            return state.channels;
        },
        currentChannel(state) {
            if (state.channels.length === 0) {
                return {};
            }
            if (state.selected == null) {
                return state.channels[0];
            }
            return state.selected || {};
        }
    }
};
