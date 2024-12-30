import express from 'express';
import { PrismaClient } from '@prisma/client';
import { authMiddleware } from '../middleware/auth';

const router = express.Router();
const prisma = new PrismaClient();

// Apply auth middleware to all carpool routes
router.use(authMiddleware);

// Get all carpools for the authenticated user
router.get('/api/carpools', async (req, res) => {
  try {
    const carpools = await prisma.carpool.findMany({
      where: {
        userId: req.user!.id // Only get carpools for this user
      }
    });
    res.json(carpools);
  } catch (error) {
    console.error('Error fetching carpools:', error);
    res.status(500).json({ error: 'Error fetching carpools' });
  }
});

// Create a new carpool for the authenticated user
router.post('/api/carpools', async (req, res) => {
  try {
    const carpool = await prisma.carpool.create({
      data: {
        ...req.body,
        userId: req.user!.id // Automatically assign to current user
      }
    });
    res.status(201).json(carpool);
  } catch (error) {
    console.error('Error creating carpool:', error);
    res.status(500).json({ error: 'Error creating carpool' });
  }
});

export default router; 