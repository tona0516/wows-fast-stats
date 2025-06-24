<script lang="ts">
  import {
    showAddAlertPlayerModal,
    showDeleteAlertPlayerModal,
    showUpdateAlertPlayerModal,
    storedAlertPlayers,
  } from "src/stores";
</script>

<div class="mb-4 flex justify-end">
  <button class="btn btn-primary" on:click={() => showAddAlertPlayerModal()}
    >追加</button
  >
</div>

<table class="table text-nowrap">
  <thead>
    <tr>
      {#each ["ID", "プレイヤー名", "アイコン", "メッセージ", "操作"] as header}
        <th class="text-center">{header}</th>
      {/each}
    </tr>
  </thead>
  <tbody>
    {#each $storedAlertPlayers as player, i}
      <tr class="hover:bg-base-200">
        <td class="text-center">{player.account_id}</td>
        <td class="text-center">{player.name}</td>
        <td class="text-center"><i class={`bi ${player.pattern}`}></i></td>
        <td class="text-center">{player.message}</td>
        <td class="text-center">
          <div class="flex justify-center">
            <button
              class="btn btn-sm btn-info mx-2"
              on:click={() =>
                showUpdateAlertPlayerModal(
                  player.account_id,
                  player.name,
                  player.pattern,
                  player.message,
                )}>編集</button
            >
            <button
              class="btn btn-sm btn-error mx-2"
              on:click={() => showDeleteAlertPlayerModal(player.account_id)}
              >削除</button
            >
          </div>
        </td>
      </tr>
    {/each}
  </tbody>
</table>
