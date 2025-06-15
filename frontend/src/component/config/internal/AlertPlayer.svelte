<script lang="ts">
import { storedAlertPlayers } from "src/stores";
import { createEventDispatcher } from "svelte";

const dispatch = createEventDispatcher();
</script>

<div>
  <h5>プレイヤー検出機能</h5>
  <ul>
    <li>戦闘情報のテーブル内のプレイヤーにアイコンを表示</li>
    <li>マウスオーバーでメモを表示</li>
    <li>プレイヤー名をクリックで追加・削除が可能</li>
  </ul>
</div>

{#if $storedAlertPlayers.length !== 0}
  <table class="table">
    <thead>
      <tr>
        {#each ["プレイヤー名", "アイコン", "メモ"] as column}
          <th>{column}</th>
        {/each}
      </tr>
    </thead>
    <tbody>
      {#each $storedAlertPlayers as player}
        <tr>
          <td>
            <details class="dropdown">
              <summary class="btn m-1">
                {player.name}
                <i class="bi bi-chevron-down"></i>
              </summary>
              <ul
                class="menu dropdown-content rounded-box z-1 w-52 p-2 shadow-sm"
              >
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
            </details>
          </td>
          <td><i class="bi {player.pattern}" /></td>
          <td>{player.message}</td>
        </tr>
      {/each}
    </tbody>
  </table>
{/if}

<div>
  <button
    class="btn btn-primary"
    on:click={() => {
      dispatch("AddAlertPlayer");
    }}>追加</button
  >
</div>
