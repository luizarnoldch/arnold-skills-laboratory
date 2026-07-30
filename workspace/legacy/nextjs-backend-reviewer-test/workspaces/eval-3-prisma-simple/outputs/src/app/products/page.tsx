import HydrateProducts from "@/features/product/hooks/HydrateProducts"
import ProductView from "@/features/product/views/ProductView"

type ProductsPageProps = {
  params: Promise<{ slug: string }>
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

const ProductsPage = async ({ params, searchParams }: ProductsPageProps) => {
  return (
    <HydrateProducts>
      <ProductView />
    </HydrateProducts>
  )
}

export default ProductsPage
