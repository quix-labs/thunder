export function injectWebhookListener() {
    const {refresh: refreshProcessors} = useProcessors()
    onMounted(() => {
        const {on} = useEvents()

        on('processor-updated', () => refreshProcessors())
        on('processor-created', () => refreshProcessors())
        on('processor-deleted', () => refreshProcessors())

        on('processor-indexed', (processorId) => {
            useToast().add({
                color: "green",
                title: `Processor ${processorId} indexed`,
            })
        })
    })

}
