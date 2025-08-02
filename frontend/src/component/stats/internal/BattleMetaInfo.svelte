<script lang="ts">
  import type { data } from "wailsjs/go/models";

  export let meta: data.Meta;

  const metaData = [
    {
      title: "マップ",
      text: meta.arena,
    },
    {
      title: "種別",
      text: meta.type,
    },
    {
      title: "取得時刻",
      text: new Date(meta.unixtime * 1000).toLocaleString(),
    },
  ];
</script>

<div class="stats shadow">
  {#each metaData as md}
    <div class="stat px-4 py-2">
      <div class="stat-title">{md.title}</div>
      <div class="stat-value text-sm">{md.text}</div>
    </div>
  {/each}
</div>
