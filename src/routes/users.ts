import express from 'express';
import { PrismaClient } from '@prisma/client';

const router = express.Router();
const prisma = new PrismaClient();

router.post('/api/users', async (req, res) => {
  try {
    const { clerk_id, email } = req.body;
    
    const user = await prisma.user.create({
      data: {
        clerk_id,
        email,
      },
    });
    
    res.status(201).json(user);
  } catch (error) {
    console.error('Error creating user:', error);
    res.status(500).json({ error: 'Error creating user' });
  }
});

export default router; 