<script lang="ts">
  import { format, fromUnixTime } from "date-fns";
  import { Card } from "sveltestrap";
  import type { domain } from "wailsjs/go/models";

  export let meta: domain.Meta;

  $: battleMetas = [
    {
      icon: "bi bi-clock-fill",
      text: formattedDate(meta.unixtime),
    },
    { icon: "bi bi-tag-fill", text: meta.type },
    { icon: "bi bi-geo-alt-fill", text: meta.arena },
  ];

  const formattedDate = (unixtime: number): string => {
    return format(fromUnixTime(unixtime), "yyyy/MM/dd HH:mm:ss");
  };
</script>

{#if battleMetas}
  <Card body class="center p-2">
    <div class="d-flex">
      {#each battleMetas as meta}
        <div class="me-2">
          <i class={meta.icon}></i>{meta.text}
        </div>
      {/each}
    </div>
  </Card>
{/if}
