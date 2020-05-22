export function conflictCustomValidity(
    component: Vue,
    type: string,
    id: string | null = null,
): (value: string) => Promise<string> {
    return (value: string): Promise<string> => {
        if (!value) {
            return Promise.resolve("");
        }

        return new Promise(resolve => {
            return component.$http
                .get(`/api/v1/${type}/${value}?excerpt=true`)
                .then(response => response.json())
                .then(
                    (response: APIResponse<ObjectMeta>) => {
                        resolve(response.data?.id !== id ? (component.$t(`messages.${type}.conflict`) as string) : "");
                    },
                    () => {
                        resolve("");
                    },
                );
        });
    };
}
