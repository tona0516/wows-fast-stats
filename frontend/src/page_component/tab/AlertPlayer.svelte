<script lang="ts">
  import { storedAlertPlayers } from "src/stores";
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();
</script>

<div class="uk-padding-small">
  <div>
    <h5>プレイヤー検出機能</h5>
    <ul class="uk-list uk-list-disc uk-list-collapse">
      <li>戦闘情報のテーブル内のプレイヤーにアイコンを表示</li>
      <li>マウスオーバーでメモを表示</li>
      <li>プレイヤー名をクリックで追加・削除が可能</li>
    </ul>
  </div>
</div>

{#if $storedAlertPlayers.length !== 0}
  <div class="uk-padding-small uk-overflow-auto">
    <table
      class="uk-table uk-table-shrink uk-table-divider uk-table-small uk-table-middle uk-text-nowrap"
    >
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
                <span uk-drop-parent-icon></span>
              </a>

              <div uk-dropdown="mode: click" uk-toggle>
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
              </div>
            </td>
            <td class="uk-text-center"><i class="bi {player.pattern}" /></td>
            <td class="uk-text-center">{player.message}</td>
          </tr>
        {/each}
      </tbody>
    </table>
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
