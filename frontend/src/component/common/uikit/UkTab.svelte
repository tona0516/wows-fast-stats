<script lang="ts">
// @ts-nocheck
export let clazz: string;
export let id: string;
</script>

<ul class={clazz} uk-tab="connect: #{id}; animation: uk-animation-fade">
  <slot />
</ul>
