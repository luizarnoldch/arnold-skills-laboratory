import HydrateProducts from "@/features/product/hooks/HydrateProducts"
import ProductView from "@/features/product/views/ProductView"

type ProductsPageProps = {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

const ProductsPage = async ({ searchParams }: ProductsPageProps) => {
  return (
    <HydrateProducts>
      <ProductView />
    </HydrateProducts>
  )
}

export default ProductsPage
