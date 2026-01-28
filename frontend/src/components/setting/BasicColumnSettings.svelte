<script lang="ts">
  import { storedDisplayPref } from "@libs/stores";
    import type { PlayerNameColorType } from "@libs/types";
  import { SaveDisplayPref } from "@wails/go/main/App";
  import { get } from "svelte/store";

  const PLAYER_NAME_COLORS: Readonly<Map<PlayerNameColorType, string>> =
  new Map<PlayerNameColorType, string>([
    ["shipPR", "艦成績のPR"],
    ["overallPR", "総合成績のPR"],
    ["threatLevel", "戦力評価"],
    ["none", "なし"],
  ]);

  const onChangePref = async () => {
    const pref = get(storedDisplayPref);
    if (!pref) return;
    await SaveDisplayPref(JSON.stringify(pref, null, 2));
  };
</script>

<div
  class="overflow-x-auto rounded-xl shadow border border-base-300 bg-base-200 mb-6"
>
  <table class="table w-full text-nowrap">
    <thead>
      <tr class="bg-base-300 text-base-content font-semibold">
        <th class="px-4 py-2">カラム名</th>
        <th class="px-4 py-2"></th>
      </tr>
    </thead>
    <tbody>
      <tr class="hover:bg-base-100">
        <td class="px-4 py-2 align-top">プレイヤー名</td>
        <td class="px-4 py-2">
          <div class="flex flex-col gap-2">
            <label class="flex items-center gap-2">
              <input
                class="toggle toggle-success"
                type="checkbox"
                bind:checked={$storedDisplayPref.player.enableNationFlag}
                on:change={onChangePref}
              />
              <span
                >クラン国籍の国旗を表示する（クラン説明から言語検出）</span
              >
            </label>
            <label class="flex items-center gap-2">
              <span>背景色タイプ</span>
              <select
                class="select select-sm select-bordered"
                bind:value={$storedDisplayPref.player.colorType}
                on:change={onChangePref}
              >
                {#each PLAYER_NAME_COLORS as color}
                  <option
                    selected={color[0] ===
                      $storedDisplayPref.player.colorType}
                    value={color[0]}>{color[1]}</option
                  >
                {/each}
              </select>
            </label>
          </div>
        </td>
      </tr>
      <tr class="hover:bg-base-100">
        <td class="px-4 py-2 align-top">艦名</td>
        <td class="px-4 py-2">
          <div class="flex flex-col gap-2">
            <label class="flex items-center gap-2">
              <input
                class="toggle toggle-success"
                type="checkbox"
                bind:checked={$storedDisplayPref.warship.enableNationFlag}
                on:change={onChangePref}
              />
              <span>国旗を表示する</span>
            </label>
          </div>
        </td>
      </tr>
    </tbody>
  </table>
</div>
