import { CardSkeleton } from "@/components/CardSkeleton";

export default function ProductLoading() {
  return (
    <div className="w-full max-w-8xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div className="flex flex-col items-center mb-12">
        <div className="h-8 w-48 bg-gray-200 rounded animate-pulse mb-4"></div>
        <div className="h-1 w-20 bg-gray-200 rounded mb-8"></div>
        <div className="flex items-center justify-between mb-6 mx-auto">
          <div className="flex gap-4 items-center">
            <div className="w-5 h-5 bg-gray-200 rounded animate-pulse"></div>
            <div className="flex gap-4 overflow-x-auto no-scrollbar w-full">
              {[...Array(5)].map((_, i) => (
                <div
                  key={i}
                  className="h-10 w-24 bg-gray-200 rounded-lg animate-pulse shrink-0"
                ></div>
              ))}
            </div>
          </div>
        </div>
      </div>

      <div className="w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
        <div className="grid grid-cols-1 sm:grid-cols-3 lg:grid-cols-4 gap-6">
          {[...Array(12)].map((_, i) => (
            <CardSkeleton key={i} />
          ))}
        </div>
      </div>
    </div>
  );
}


