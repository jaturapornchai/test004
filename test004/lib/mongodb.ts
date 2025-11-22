// ========================================
// MongoDB Atlas API Service
// เชื่อมต่อกับ API Endpoint ที่กำหนด
// ========================================

const API_BASE_URL = 'https://smlgoapi.dedepos.com/v1';

// ========== Types ==========
interface MongoFilter {
  [key: string]: unknown;
}

interface MongoProjection {
  [key: string]: 0 | 1;
}

interface MongoSort {
  [key: string]: 1 | -1;
}

interface UpdateOptions {
  upsert?: boolean;
  replaceone?: boolean;
}

interface DeleteOptions {
  delete_many?: boolean;
}

interface GetOptions {
  projection?: MongoProjection;
  sort?: MongoSort;
  limit?: number;
  skip?: number;
}

// ========== API Functions ==========

/**
 * Insert หรือ Update document (Upsert)
 */
export async function mongoUpdate<T>(
  collection: string,
  filter: MongoFilter,
  data: Partial<T>,
  options: UpdateOptions = {}
): Promise<{
  status: string;
  code: number;
  matched_count: number;
  modified_count: number;
  upserted_id: string | null;
}> {
  const response = await fetch(`${API_BASE_URL}/mongoatlasupdate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      collection,
      filter,
      data,
      upsert: options.upsert ?? true,
      replaceone: options.replaceone ?? false,
    }),
  });

  if (!response.ok) {
    throw new Error(`MongoDB update failed: ${response.statusText}`);
  }

  return response.json();
}

/**
 * Delete document(s)
 */
export async function mongoDelete(
  collection: string,
  filter: MongoFilter,
  options: DeleteOptions = {}
): Promise<{
  status: string;
  code: number;
  deleted_count: number;
}> {
  const response = await fetch(`${API_BASE_URL}/mongoatlasdelete`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      collection,
      filter,
      delete_many: options.delete_many ?? false,
    }),
  });

  if (!response.ok) {
    throw new Error(`MongoDB delete failed: ${response.statusText}`);
  }

  return response.json();
}

/**
 * Query documents
 */
export async function mongoGet<T>(
  collection: string,
  filter: MongoFilter = {},
  options: GetOptions = {}
): Promise<{
  status: string;
  code: number;
  count: number;
  data: T[];
}> {
  const response = await fetch(`${API_BASE_URL}/mongoatlasget`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      collection,
      filter,
      projection: options.projection,
      sort: options.sort,
      limit: options.limit,
      skip: options.skip,
    }),
  });

  if (!response.ok) {
    throw new Error(`MongoDB get failed: ${response.statusText}`);
  }

  return response.json();
}

// ========== Helper Functions ==========

/**
 * สร้าง ID ใหม่
 */
export function generateId(prefix: string = ''): string {
  const timestamp = Date.now().toString(36);
  const random = Math.random().toString(36).substring(2, 8);
  return prefix ? `${prefix}_${timestamp}${random}` : `${timestamp}${random}`;
}

/**
 * สร้าง Filter สำหรับ Multi-tenant
 */
export function withCoopId(coopId: string, filter: MongoFilter = {}): MongoFilter {
  return {
    ...filter,
    coopId,
  };
}

/**
 * สร้าง Pagination
 */
export function getPagination(page: number, limit: number): { skip: number; limit: number } {
  return {
    skip: (page - 1) * limit,
    limit,
  };
}

// ========== CRUD Wrapper Functions ==========

/**
 * Generic CRUD operations with multi-tenant support
 */
export class MongoCollection<T extends { coopId: string }> {
  constructor(private collectionName: string) {}

  // Create
  async create(data: T): Promise<T> {
    const now = new Date().toISOString();
    const docWithTimestamp = {
      ...data,
      createdAt: now,
      updatedAt: now,
    };

    await mongoUpdate(
      this.collectionName,
      { coopId: data.coopId, _id: (data as T & { _id?: string })._id },
      docWithTimestamp,
      { upsert: true }
    );

    return docWithTimestamp as T;
  }

  // Read one
  async findOne(coopId: string, filter: MongoFilter): Promise<T | null> {
    const result = await mongoGet<T>(
      this.collectionName,
      withCoopId(coopId, filter),
      { limit: 1 }
    );

    return result.data[0] || null;
  }

  // Read many
  async find(
    coopId: string,
    filter: MongoFilter = {},
    options: GetOptions = {}
  ): Promise<{ data: T[]; count: number }> {
    const result = await mongoGet<T>(
      this.collectionName,
      withCoopId(coopId, filter),
      options
    );

    return {
      data: result.data,
      count: result.count,
    };
  }

  // Update
  async update(
    coopId: string,
    filter: MongoFilter,
    data: Partial<T>
  ): Promise<boolean> {
    const result = await mongoUpdate(
      this.collectionName,
      withCoopId(coopId, filter),
      {
        ...data,
        updatedAt: new Date().toISOString(),
      },
      { upsert: false }
    );

    return result.modified_count > 0;
  }

  // Delete
  async delete(coopId: string, filter: MongoFilter): Promise<boolean> {
    const result = await mongoDelete(
      this.collectionName,
      withCoopId(coopId, filter)
    );

    return result.deleted_count > 0;
  }

  // Delete many
  async deleteMany(coopId: string, filter: MongoFilter): Promise<number> {
    const result = await mongoDelete(
      this.collectionName,
      withCoopId(coopId, filter),
      { delete_many: true }
    );

    return result.deleted_count;
  }
}

// ========== Collection Instances ==========

import type {
  Cooperative,
  Member,
  ShareAccount,
  DepositAccount,
  DepositTransaction,
  LoanProduct,
  LoanApplication,
  LoanContract,
  LoanPayment,
  Product,
  StockMovement,
  Sale,
  ProduceCollection,
  ChartOfAccount,
  JournalEntry,
  ProfitAllocation,
  DividendDistribution,
  User,
} from './types';

// สหกรณ์
export const cooperatives = new MongoCollection<Cooperative>('cooperatives');

// สมาชิก
export const members = new MongoCollection<Member>('members');

// ทุนเรือนหุ้น
export const shareAccounts = new MongoCollection<ShareAccount>('share_accounts');

// เงินฝาก
export const depositAccounts = new MongoCollection<DepositAccount>('deposit_accounts');
export const depositTransactions = new MongoCollection<DepositTransaction>('deposit_transactions');

// สินเชื่อ
export const loanProducts = new MongoCollection<LoanProduct>('loan_products');
export const loanApplications = new MongoCollection<LoanApplication>('loan_applications');
export const loanContracts = new MongoCollection<LoanContract>('loan_contracts');
export const loanPayments = new MongoCollection<LoanPayment>('loan_payments');

// ธุรกิจ
export const products = new MongoCollection<Product>('products');
export const stockMovements = new MongoCollection<StockMovement>('stock_movements');
export const sales = new MongoCollection<Sale>('sales');
export const produceCollections = new MongoCollection<ProduceCollection>('produce_collections');

// บัญชี
export const chartOfAccounts = new MongoCollection<ChartOfAccount>('chart_of_accounts');
export const journalEntries = new MongoCollection<JournalEntry>('journal_entries');
export const profitAllocations = new MongoCollection<ProfitAllocation>('profit_allocations');
export const dividendDistributions = new MongoCollection<DividendDistribution>('dividend_distributions');

// ผู้ใช้งาน
export const users = new MongoCollection<User>('users');
