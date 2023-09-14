import type { KeyValue } from "src/lib/value_object/KeyValue";

export class ArrayMap<K, V> extends Map<K, V> {
  toArray(): KeyValue<K, V>[] {
    const result: KeyValue<K, V>[] = [];
    this.forEach((value, key) => {
      result.push({ key: key, value: value });
    });
    return result;
  }
}
