export function mapReferences(refs: Array<Reference>): Record<string, unknown> {
    return (
        refs.reduce((out: Record<string, unknown>, ref: Reference) => {
            switch (ref.type) {
                case "chart":
                    out[`chart|${(ref.value as Chart).id}`] = ref.value;
                    break;
            }

            return out;
        }, {}) ?? {}
    );
}
