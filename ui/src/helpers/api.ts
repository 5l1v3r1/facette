export function conflictCustomValidity(component: Vue, type: string): (value: string) => Promise<string> {
    return (value: string): Promise<string> => {
        if (!value) {
            return Promise.resolve("");
        }

        return new Promise(resolve => {
            return component
                .$http({
                    url: `/api/v1/${type}/${value}`,
                    method: "HEAD",
                })
                .then(
                    () => {
                        resolve(component.$t(`messages.${type}.conflict`) as string);
                    },
                    () => {
                        resolve("");
                    },
                );
        });
    };
}
