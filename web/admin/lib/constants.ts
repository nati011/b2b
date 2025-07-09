export const ACL: Record<string, string[]> = {
    '/': ['order', 'order_distributors'],
    '/distributors':['distributor'],
    '/distributors/agents':['distributor_user'],
    '/distributor/agents/form':['create_distributor_user'],
    '/invoice':['invoice'],
    '/orders':['order'],
    '/products':['product'],
    '/retailers':['retailer'],
    '/role':['role'],
    '/transactions':['transaction']
}
