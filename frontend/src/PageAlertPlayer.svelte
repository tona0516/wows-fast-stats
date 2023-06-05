<script lang="ts">
import { createEventDispatcher } from "svelte";
import type { vo } from "wailsjs/go/models";
import clone from "clone";
import { alertPlayers } from "./stores.js";

const dispatch = createEventDispatcher();

let players: vo.AlertPlayer[];
alertPlayers.subscribe((it) => (players = it));
</script>

<div class="mt-3 alert-player">
  <!-- alert players -->
  <div class="m-2">
    {#if players.length === 0}
      <p>プレイヤーリストがありません</p>
    {:else}
      <div class="d-flex flex-row centerize">
        <table class="table table-sm table-text-color w-auto">
          <thead>
            <tr>
              <th>プレイヤー名</th>
              <th>アイコン</th>
              <th>メモ</th>
              <th></th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {#each players as player}
              <tr>
                <td>{player.name}</td>
                <td><i class="bi {player.pattern}"></i></td>
                <td>{player.message}</td>
                <td>
                  <button
                    type="button"
                    class="my-1 btn btn-sm btn-success"
                    on:click="{() => {
                      const target = clone(player);
                      dispatch('UpdateAlertPlayer', { target: target });
                    }}">編集</button
                  >
                </td>
                <td>
                  <button
                    type="button"
                    class="my-1 btn btn-sm btn-danger"
                    on:click="{() => {
                      const target = clone(player);
                      dispatch('RemoveAlertPlayer', { target: target });
                    }}">削除</button
                  >
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </div>

  <!-- add button -->
  <button
    type="button"
    class="btn btn-sm btn-primary"
    on:click="{() => dispatch('AddAlertPlayer')}"
  >
    追加
  </button>
</div>
