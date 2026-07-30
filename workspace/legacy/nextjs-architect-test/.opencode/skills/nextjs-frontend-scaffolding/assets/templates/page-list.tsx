import Hydrate[Entity]s from "@/features/[entity-kebab]/hooks/Hydrate[Entity]s"
import [Entity]View from "@/features/[entity-kebab]/views/[Entity]View"

type [Entity]sPageProps = {
  params: Promise<{ [slugId]: string }>
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

const [Entity]sPage = async ({ params, searchParams }: [Entity]sPageProps) => {
  // const { slugId } = await params <- Uncomment if required
  // const value = (await searchParams).key <- Uncomment if required

  return (
    <Hydrate[Entity]s>
      <[Entity]View />
    </Hydrate[Entity]s>
  )
}

export default [Entity]sPage
