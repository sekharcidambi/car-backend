import { Request, Response, NextFunction } from 'express';
import { PrismaClient } from '@prisma/client';
import jwt from 'jsonwebtoken';

const prisma = new PrismaClient();

declare global {
  namespace Express {
    interface Request {
      user?: {
        id: string;
        clerk_id: string;
      };
    }
  }
}

export async function authMiddleware(req: Request, res: Response, next: NextFunction) {
  try {
    const token = req.headers.authorization?.split(' ')[1];
    
    if (!token) {
      return res.status(401).json({ error: 'No token provided' });
    }

    // Verify the JWT token from Clerk
    const decoded = jwt.verify(token, process.env.CLERK_JWT_PUBLIC_KEY!, {
      algorithms: ['RS256'],
    });

    // Get clerk_id from the decoded token
    const clerk_id = (decoded as any).sub;

    // Find the user in your database
    const user = await prisma.user.findUnique({
      where: { clerk_id },
    });

    if (!user) {
      return res.status(401).json({ error: 'User not found' });
    }

    // Attach user to request object
    req.user = user;
    next();
  } catch (error) {
    console.error('Auth error:', error);
    res.status(401).json({ error: 'Invalid token' });
  }
} 