<div class="uk-overflow-auto">
  <table
    class="uk-table uk-table-shrink uk-table-divider uk-table-small uk-table-middle uk-text-nowrap uk-margin-bottom-small"
  >
    <slot />
  </table>
</div>

<style>
  :global(.uk-table td) {
    padding-top: 2px;
    padding-bottom: 2px;
    padding-left: 4px;
    padding-right: 4px;
  }

  :global(.uk-table th) {
    padding-top: 2px;
    padding-bottom: 2px;
    padding-left: 4px;
    padding-right: 4px;
  }
</style>
