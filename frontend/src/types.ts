export interface Account {
    id: string,
    name: string,
    icon: string
}

export interface AccountPatch{
    name?: string,
    icon?: string
}