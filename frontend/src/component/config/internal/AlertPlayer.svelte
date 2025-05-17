<script lang="ts">
import UkDowndown from "src/component/common/uikit/UkDowndown.svelte";
import UkIcon from "src/component/common/uikit/UkIcon.svelte";
import UkTable from "src/component/common/uikit/UkTable.svelte";
import { storedAlertPlayers } from "src/stores";
import { createEventDispatcher } from "svelte";

const dispatch = createEventDispatcher();
</script>

<div class="uk-padding-small">
  <h5>プレイヤー検出機能</h5>
  <ul class="uk-list uk-list-disc uk-list-collapse">
    <li>戦闘情報のテーブル内のプレイヤーにアイコンを表示</li>
    <li>マウスオーバーでメモを表示</li>
    <li>プレイヤー名をクリックで追加・削除が可能</li>
  </ul>
</div>

{#if $storedAlertPlayers.length !== 0}
  <div class="uk-padding-small">
    <UkTable>
      <thead>
        <tr>
          {#each ["プレイヤー名", "アイコン", "メモ"] as column}
            <th class="uk-text-center">{column}</th>
          {/each}
        </tr>
      </thead>
      <tbody>
        {#each $storedAlertPlayers as player}
          <tr>
            <td class="uk-text-center">
              <!-- svelte-ignore a11y-invalid-attribute -->
              <a href="#">
                {player.name}
                <UkIcon name="chevron-down" />
              </a>

              <UkDowndown>
                <ul class="uk-nav uk-dropdown-nav">
                  <li>
                    <!-- svelte-ignore a11y-invalid-attribute -->
                    <a
                      href="#"
                      on:click={() => {
                        dispatch("EditAlertPlayer", { target: player });
                      }}>編集</a
                    >
                  </li>
                  <li>
                    <!-- svelte-ignore a11y-invalid-attribute -->
                    <a
                      href="#"
                      on:click={() => {
                        dispatch("RemoveAlertPlayer", { target: player });
                      }}>削除</a
                    >
                  </li>
                </ul>
              </UkDowndown>
            </td>
            <td class="uk-text-center"><i class="bi {player.pattern}" /></td>
            <td class="uk-text-center">{player.message}</td>
          </tr>
        {/each}
      </tbody>
    </UkTable>
  </div>
{/if}

<div class="uk-padding-small">
  <button
    class="uk-button uk-button-primary uk-text-nowrap"
    on:click={() => {
      dispatch("AddAlertPlayer");
    }}>追加</button
  >
</div>
