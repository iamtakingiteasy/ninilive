<template>
  <div id="root">
    <div id="wrapper">
      <router-view/>
    </div>
  </div>
</template>

<style>
  @font-face{
    font-family: 'Livechan';
    src: url('~@/assets/Roboto-Medium.ttf');
  }
  @font-face{
    font-family: 'Livechan Mono';
    src: url('~@/assets/Consolas.ttf');
  }
  #root {
    font-family: Livechan, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    display: flex;
    flex-direction: row;
    width: 100%;
    height: 100%;
    justify-content: center;
    line-height: 1.15;
    color: var(--fg-text-color);
    background: var(--bg-base-color);
  }
  #wrapper {
    display: flex;
    flex-direction: column;
    min-width: 10%;
    max-width: 90%;
    flex-grow: 1;
  }
  @media (max-aspect-ratio: 1/1) {
    #wrapper {
      max-width: 90%;
    }
  }
  @media (min-aspect-ratio: 1/1) {
    #wrapper {
      max-width: 70%;
    }
  }
  @media (min-aspect-ratio: 4/3) {
    #wrapper {
      max-width: 60%;
    }
  }
  @media (min-aspect-ratio: 16/9) {
    #wrapper {
      max-width: 50%;
    }
  }
  html, body {
    margin: 0;
    padding: 0;
    width: 100%;
    height: 100%;
  }
  * {
    outline: none;
  }
  :root {
    --bg-base-color: #36393F;

    --bg-panel-light: rgba(0, 0, 0, 0.05);
    --bg-panel-color: rgba(0, 0, 0, 0.3);
    --bg-panel-hover: rgba(0, 0, 0, 0.1);

    --bg-form-color: #40444B;

    --bg-focus-color: #2F3136;

    --bg-accent-color: rgba(0, 0, 0, 0.6);
    --bg-accent-hover: rgba(0, 0, 0, 0.75);

    --bg-badge-color: #D9D9EE;
    --fg-badge-color: #333333;

    --fg-text-color: #CCD7D9;
    --fg-link-color: #05B0E3;
    --fg-name-color: #31BE71;
    --fg-mod-color: #FF8A11;
    --fg-greentext-color: #789922;
  }
  a {
    color: var(--fg-link-color);
    cursor: pointer;
    text-decoration: none;
  }

  ::-webkit-scrollbar {
    width: 1ex;
  }

  ::-webkit-scrollbar-track {
    background-color: var(--bg-panel-color);
    border-radius: 1ex;
  }

  ::-webkit-scrollbar-thumb {
    background: var(--bg-accent-color);
    border-radius: 1ex;
  }

  ::-webkit-scrollbar-thumb:hover {
    background: var(--bg-accent-hover);
  }

  #form {
    width: 900px;
    border-radius: 0 0 1ex 1ex;
  }

  ::-webkit-file-upload-button {
    background: var(--bg-accent-color);
    color: #fff;
    border: 0 solid;
    border-radius: 0.5ex;
    padding: 0.4em;
    outline: none;
  }

  ::-webkit-file-upload-button:hover {
    background: var(--bg-accent-hover);
  }
</style>