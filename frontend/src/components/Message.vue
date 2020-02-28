<template>
    <div class="message">
        <div class="meta">
            <span class="name">{{ message.name || message.user.name }}</span>
            <span class="trip" v-if="message.trip">#{{ message.trip }}</span>
            <span class="badge">{{ channelName }}</span>
            <span class="time">{{ message.time }}</span>
            <a href="#" v-on:click="$emit('insertNo', message.id)">No.{{message.id}}</a>
        </div>
        <div class="body" v-html="bodyHTML" ref="body"></div>
        <div class="file" v-if="message.file">
            <a :href="fileURL" target="_blank">
                <img v-on:load="$emit('loaded')" :src="fileURL" :alt="message.file.name"/>
            </a>
        </div>
    </div>
</template>
<script>
    import {endpointHTTP} from "../config";

    export default {
        name: 'Message',
        props: {
            message: Object,
            channelName: String
        },
        mounted() {
            this.$refs.body.querySelectorAll('.topost').forEach(x => {
                x.onclick = () => this.$emit('scrollToPost', x.dataset.id);
            });
        },
        computed: {
            endpointHTTP() {
                return endpointHTTP;
            },
            bodyHTML() {
                const text = document.createTextNode(this.message.body);
                const p = document.createElement('p');
                p.appendChild(text);
                let raw = p.innerHTML;
                raw = raw.replace(/&gt;&gt;([0-9]+)/,'<a href="#" class="topost" data-id="$1">&gt;&gt;$1</a>');
                return raw;
            },
            fileURL() {
                if (!this.message.file) {
                    return '';
                }
                return endpointHTTP + '/ws/download/' + this.message.file.path + '/' + this.message.file.name;
            }
        }
    }
</script>
<style>
    .message {
        margin: 0 1em;
        padding: 1em;
        text-align: left;
        border-top: 0.1ex solid var(--bg-accent-color);
    }
    .message > .meta > * {
        display: inline-block;
        margin-right: 1ex;
    }
    .message > .meta > .name {
        color: var(--fg-name-color);
    }
    .message > .meta > .trip {
        margin-left: -1ex;
    }
    .message > .meta > .badge {
        padding: 0 0.5ex;
        background: var(--bg-badge-color);
        color: var(--fg-badge-color);
        font-weight: bold;
        text-transform: capitalize;
    }
    .message > .body {
        white-space: pre-line;
        margin-left: 1em;
    }
    .message > .file > a > img {
        margin-top: 1em;
        max-width: 100%;
    }
</style>