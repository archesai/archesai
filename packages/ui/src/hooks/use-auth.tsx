export const useAuth = () => {
  const signIn = () => {
    localStorage.setItem('isAuthenticated', 'true')
  }

  const signOut = () => {
    localStorage.removeItem('isAuthenticated')
  }

  const isLogged = () => localStorage.getItem('isAuthenticated') === 'true'

  return { isLogged, signIn, signOut }
}

export type AuthContext = ReturnType<typeof useAuth>
