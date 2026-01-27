<script lang="ts">
  import { STATS_KEYS, STATS_COLUMN_INFO } from "@libs/constants";
  import { storedDisplayPref } from "@libs/stores";
  import { SaveDisplayPref } from "@wails/go/main/App";
  import { get } from "svelte/store";

  const onChangePref = async () => {
    const pref = get(storedDisplayPref);
    if (!pref) return;
    await SaveDisplayPref(JSON.stringify(pref, null, 2));
  };
</script>

<div
  class="overflow-x-auto rounded-xl shadow border border-base-300 bg-base-200"
>
  <table class="table w-full text-nowrap">
    <thead>
      <tr class="bg-base-300 text-base-content font-semibold">
        {#each ["カラム名", "艦成績", "総合成績", "小数点以下の桁数"] as columns}
          <th class="text-center px-4 py-2">{columns}</th>
        {/each}
      </tr>
    </thead>
    <tbody>
      {#each STATS_KEYS as statsKey}
        {@const info = STATS_COLUMN_INFO[statsKey]}
        <tr class="hover:bg-base-100">
          <td class="px-4 py-2">{info.fullName}</td>
          {#if "showShip" in $storedDisplayPref[statsKey] && typeof $storedDisplayPref[statsKey].showShip === 'boolean'}
            <td class="text-center px-4 py-2">
              <input
                class="toggle toggle-success"
                type="checkbox"
                bind:checked={$storedDisplayPref[statsKey].showShip}
                on:change={onChangePref}
              />
            </td>
          {:else}
            <td></td>
          {/if}
          {#if "showOverall" in $storedDisplayPref[statsKey] && typeof $storedDisplayPref[statsKey].showOverall === 'boolean'}
            <td class="text-center px-4 py-2">
              <input
                class="toggle toggle-success"
                type="checkbox"
                bind:checked={
                  $storedDisplayPref[statsKey].showOverall
                }
                on:change={onChangePref}
              />
            </td>
          {:else}
            <td></td>
          {/if}
          {#if "digit" in $storedDisplayPref[statsKey] && typeof $storedDisplayPref[statsKey].digit === 'number'}
            <td class="text-center px-4 py-2">
              <select
                class="select select-sm select-bordered"
                bind:value={$storedDisplayPref[statsKey].digit}
                on:change={onChangePref}
              >
                {#each [0, 1, 2] as digit}
                  <option
                    selected={digit ===
                      $storedDisplayPref[statsKey].digit}
                    value={digit}>{digit}</option
                  >
                {/each}
              </select>
            </td>
          {:else}
            <td></td>
          {/if}
        </tr>
      {/each}
    </tbody>
  </table>
</div>
