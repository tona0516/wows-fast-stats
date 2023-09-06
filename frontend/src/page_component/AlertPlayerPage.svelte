<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import clone from "clone";
  import { storedAlertPlayers, storedUserConfig } from "src/stores";
  import type { domain } from "wailsjs/go/models";
  import { Button } from "sveltestrap";

  const dispatch = createEventDispatcher();

  const alertPlayerColumns = ["プレイヤー名", "アイコン", "メモ", "操作"];

  function onClickAdd() {
    dispatch("AddAlertPlayer");
  }

  function onClickEdit(player: domain.AlertPlayer) {
    const target = clone(player);
    dispatch("UpdateAlertPlayer", { target: target });
  }

  function onClickRemove(player: domain.AlertPlayer) {
    const target = clone(player);
    dispatch("RemoveAlertPlayer", { target: target });
  }
</script>

<div class="mt-3 center">
  <!-- introduction -->
  <div class="alert alert-primary">
    <p>
      <i class="bi bi-info-circle-fill" /> プレイヤー検出機能
    </p>
    <ul class="m-0">
      <li>リストに追加されたプレイヤーにアイコン表示</li>
      <li>マッチのプレイヤー名クリックからも追加・削除可能</li>
    </ul>
  </div>

  <!-- alert players -->
  <div class="m-2">
    {#if $storedAlertPlayers.length === 0}
      <p>プレイヤーリストがありません</p>
    {:else}
      <table class="table table-sm table-text-color w-auto td-multiple">
        <thead>
          <tr>
            {#each alertPlayerColumns as column}
              <th>{column}</th>
            {/each}
          </tr>
        </thead>
        <tbody>
          {#each $storedAlertPlayers as player}
            <tr>
              <td>{player.name}</td>
              <td><i class="bi {player.pattern}" /></td>
              <td>{player.message}</td>
              <td>
                <Button
                  size="sm"
                  color="success"
                  on:click={() => onClickEdit(player)}
                  >編集
                </Button>
                <Button
                  size="sm"
                  color="danger"
                  on:click={() => onClickRemove(player)}>削除</Button
                >
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    {/if}
  </div>

  <!-- add button -->
  <Button size="sm" color="primary" on:click={onClickAdd}>追加</Button>
</div>
