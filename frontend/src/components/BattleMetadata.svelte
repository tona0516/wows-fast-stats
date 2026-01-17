<script lang="ts">
  import { core } from "@wails/go/models";

  export let metadata: core.BattleMetaData;

  const metaTexts = [
    {
      title: "マップ",
      text: metadata.arena,
    },
    {
      title: "種別",
      text: metadata.type,
    },
    {
      title: "取得時刻",
      text: new Date(metadata.unixtime * 1000).toLocaleString(),
    },
  ];
</script>

<div class="stats shadow">
  {#each metaTexts as md}
    <div class="stat px-4 py-2">
      <div class="stat-title">{md.title}</div>
      <div class="stat-value text-sm">{md.text}</div>
    </div>
  {/each}
</div>
