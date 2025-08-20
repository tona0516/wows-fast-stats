<script lang="ts">
  export let colspan: number;
  export let text: string;
</script>

<td class="p-1 text-center" {colspan}>{text}</td>
