<script lang="ts">
import { createEventDispatcher } from "svelte";
import clone from "clone";
import { storedAlertPlayers, storedUserConfig } from "../stores.js";
import { get } from "svelte/store";
import type { vo } from "../../wailsjs/go/models.js";
import { Button } from "sveltestrap";

const dispatch = createEventDispatcher();

let alertPlayers = get(storedAlertPlayers);
storedAlertPlayers.subscribe((it) => (alertPlayers = it));

let userConfig = get(storedUserConfig);
storedUserConfig.subscribe((it) => {
  userConfig = it;
});

const alertPlayerColumns = ["プレイヤー名", "アイコン", "メモ", "操作"];

function onClickAdd() {
  dispatch("AddAlertPlayer");
}

function onClickEdit(player: vo.AlertPlayer) {
  const target = clone(player);
  dispatch("UpdateAlertPlayer", { target: target });
}

function onClickRemove(player: vo.AlertPlayer) {
  const target = clone(player);
  dispatch("RemoveAlertPlayer", { target: target });
}
</script>

<div class="mt-3 center">
  <!-- introduction -->
  <div class="alert alert-primary">
    <p>
      <i class="bi bi-info-circle-fill"></i> プレイヤー検出機能
    </p>
    <ul class="m-0">
      <li>リストに追加されたプレイヤーにアイコン表示</li>
      <li>マッチのプレイヤー名クリックからも追加・削除可能</li>
    </ul>
  </div>

  <!-- alert players -->
  <div class="m-2">
    {#if alertPlayers.length === 0}
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
          {#each alertPlayers as player}
            <tr>
              <td>{player.name}</td>
              <td><i class="bi {player.pattern}"></i></td>
              <td>{player.message}</td>
              <td>
                <Button
                  class="m-1"
                  size="sm"
                  color="success"
                  style="font-size: {userConfig.font_size};"
                  on:click="{() => onClickEdit(player)}">編集</Button
                >
                <Button
                  class="m-1"
                  size="sm"
                  color="danger"
                  style="font-size: {userConfig.font_size};"
                  on:click="{() => onClickRemove(player)}">削除</Button
                >
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    {/if}
  </div>

  <!-- add button -->
  <Button
    class="m-1"
    size="sm"
    color="primary"
    style="font-size: {userConfig.font_size};"
    on:click="{onClickAdd}">追加</Button
  >
</div>
