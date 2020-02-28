export const usersModule = {
    namespaced: true,
    state: {
        user: null,
        name: null,
        trip: null,
        session: null,
        known: [],
    },
    mutations: {
        self(state, {user, session = {}}) {
            state.user = user;
            state.session = session;
        },
        set(state, {sessions = []}) {
            state.known = sessions.sort((a, b) => a.id.localeCompare(b.id))
        },
        add(state, {session = {}}) {
            if (!state.known.some(x => x.id === session.id)) {
                state.known = [...state.known, session].sort((a, b) => a.id.localeCompare(b.id))
            }
        },
        del(state, {session = {}}) {
            const idx = state.known.findIndex(x => x.id === session.id);
            state.known.splice(idx, 1);
        },
        setSelfUsername(state, {username}) {
            state.name = username
        },
        setSelfTripcode(state, {tripcode}) {
            state.trip = tripcode
        },
    },
    getters: {
        selfUsername(state) {
            return state.name || (state.user ? state.user.name : null) || 'Anonymous';
        },
        selfTripcode(state) {
            return state.trip || '';
        }
    }
};