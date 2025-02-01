<script lang="ts">
  import { format, fromUnixTime } from "date-fns";
  import UkIcon from "src/component/common/uikit/UkIcon.svelte";
  import type { model } from "wailsjs/go/models";

  export let meta: model.Meta;

  $: battleMetas = [
    {
      icon: "clock",
      text: format(fromUnixTime(meta.unixtime), "yyyy/MM/dd HH:mm:ss"),
    },
    {
      icon: "tag",
      text: meta.type,
    },
    {
      icon: "location",
      text: meta.arena,
    },
  ];
</script>

<div class="uk-flex uk-flex-center">
  {#each battleMetas as meta}
    <div class="uk-margin-small-right uk-text-nowrap">
      <UkIcon name={meta.icon} />
      <span class="uk-text-middle">{meta.text}</span>
    </div>
  {/each}
</div>
