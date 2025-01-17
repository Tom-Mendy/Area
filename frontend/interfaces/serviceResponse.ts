export interface User {
  username: string;
  email: string;
  id: number;
}

export interface Service {
  name: string;
  oauth: boolean;
}

export interface Token {
  service: Service;
  id: number;
}

export interface ServiceResponse {
  tokens: Token[];
  user: User;
}
