export type UserData = {
  first_name: string;
  last_name: string;
  work_email: string;
  employment_number: number | null;
  department: string;
  mailing_address: string;
  phone_number: number | null;
  password?: string;
  postal_code: string;
  city: string;
  state: string;
  country: string;
  foapa_details: FoapaData[];
};

export type UserDataPreVerification = {
  first_name: string;
  last_name: string;
  work_email: string;
  employment_number: number | null;
  phone_number: number | null;
  department: string;
};


export type FoapaData = {
  fund: string;
  organization: string;
  account: string;
  program: string;
  activity: string;
  foapa_name: string;
  initial_amount: string;
  current_amount: string;
};

export type AddressDetails = {
  name: string;
  code: string;
};

export type Activity = {
  activityName: string;
  cost: number;
  additionalInformation?: string;
  foapaNumber: string;
  activityDate: string;
  activityReceipt: string;
  activityId: string;
  _id?: string;
};

export type FoapaNumbers = {
  employmentNumber: number;
  foapaNumber: string;
  foapaName: string;
  currentAmount: number | "N/A";
  initialAmount: number | "N/A";
};

export type ReimbursementTicket = {
  reimbursementName: String;
  reimbursementReason: String;
  destination: String;
  paymentRetrievalMethod: "Hold for Pickup" | "Direct Deposit" | "";
  UDMPUVoucher: Boolean;
  totalCost: number;
  reimbursementReceipts: { url: String; id: String }[];
  reimbursementStatus: string;
  reimbursementDate: string;
  activities: Activity[];
};
