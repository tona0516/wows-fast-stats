<script lang="ts">
import { format, fromUnixTime } from "date-fns";
import type { data } from "wailsjs/go/models";

export let meta: data.Meta;

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

<div class="flex">
  {#each battleMetas as meta}
    <i class="bi bi-{meta.icon}">{meta.text}</i>
  {/each}
</div>
