<script lang="ts">
export let id: string;
export let title: string;
export let component;
</script>

<div class="accordion m-1" id="{id}">
  <div class="accordion-item">
    <h2 class="accordion-header" id="heading{id}">
      <button
        class="accordion-button"
        type="button"
        data-bs-toggle="collapse"
        data-bs-target="#collapse{id}"
        aria-expanded="true"
        aria-controls="collapse{id}"
      >
        {title}
      </button>
    </h2>
    <div
      id="collapse{id}"
      class="accordion-collapse collapse"
      aria-labelledby="heading{id}"
      data-bs-parent="#{id}"
    >
      <div class="accordion-body">
        <svelte:component this="{component}" />
      </div>
    </div>
  </div>
</div>
