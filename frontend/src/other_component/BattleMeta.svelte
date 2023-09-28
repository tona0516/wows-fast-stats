<script lang="ts">
  import { format, fromUnixTime } from "date-fns";
  import type { domain } from "wailsjs/go/models";

  export let meta: domain.Meta;

  $: battleMetas = [
    {
      icon: "clock",
      text: formattedDate(meta.unixtime),
    },
    { icon: "tag", text: meta.type },
    { icon: "location", text: meta.arena },
  ];

  const formattedDate = (unixtime: number): string => {
    return format(fromUnixTime(unixtime), "yyyy/MM/dd HH:mm:ss");
  };
</script>

<div class="uk-flex uk-flex-center">
  {#each battleMetas as meta}
    <div class="uk-margin-small-right uk-text-nowrap">
      <span uk-icon="icon: {meta.icon}" />
      <span class="uk-text-middle">{meta.text}</span>
    </div>
  {/each}
</div>
